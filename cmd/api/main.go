package main

import (
	"context"
	"log"
	"os"

	"github.com/JuanSposada/me-transfer/internal/api/handlers" // Tu ruta de handlers
	"github.com/JuanSposada/me-transfer/internal/api/middleware"
	"github.com/JuanSposada/me-transfer/internal/repository/postgres"
	"github.com/JuanSposada/me-transfer/internal/storage"
	"github.com/JuanSposada/me-transfer/internal/worker"
	"github.com/gin-gonic/gin" // Importante: go get github.com/gin-gonic/gin
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar entorno
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ Nota: Usando variables de entorno del sistema")
	}

	// 2. Inicializar DB (Postgres)
	connStr := os.Getenv("POSTGRES_URL")
	repo, err := postgres.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatalf("❌ Error crítico en DB: %v", err)
	}
	defer repo.Pool.Close()

	// 3. Inicializar Storage (Supabase)
	storageSvc := storage.NewSupabaseStorage(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		os.Getenv("SUPABASE_BUCKET"),
	)

	// 4. Verificación de salud
	if err := repo.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ La DB no responde: %v", err)
	}

	log.Println("✅ INFRAESTRUCTURA LISTA: Postgres y Supabase conectados.")

	// --- 🚀 AQUÍ ENTRA GIN (PERSONA C) ---

	// 5. Configurar el servidor Gin
	r := gin.Default()

	// Middleware de CORS básico
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 6. Inicializar el Handler
	fileHandler := handlers.NewFileHandler(repo, storageSvc)

	// 7. Definir las rutas
	r.POST("/upload", middleware.MaxAllowedSize(5*1024*1024), fileHandler.Upload)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "¡Servidor vivo!"})
	})

	// En tu main.go, dentro de la función main:
	r.GET("/download/:id", fileHandler.Download)

	cleanWorker := worker.NewCleanupWorker(repo, storageSvc)
	cleanWorker.Start(context.Background()) // Esto arranca el bucle en segundo plano
	// 8. Arrancar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Local fallback
	}
	log.Println("🚀 API corriendo en http://localhost:" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ No se pudo arrancar el servidor: %v", err)
	}

}
