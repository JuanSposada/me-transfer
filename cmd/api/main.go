package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/JuanSposada/me-transfer/internal/api/handlers"
	"github.com/JuanSposada/me-transfer/internal/repository/postgres"
	"github.com/JuanSposada/me-transfer/internal/service"
	"github.com/JuanSposada/me-transfer/internal/storage"
)

func main() {
	// 1. Cargar entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Usando variables de entorno del sistema")
	}

	// 2. Inicializar DB (Postgres)
	connStr := os.Getenv("POSTGRES_URL")
	repo, err := postgres.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatalf("❌ Error DB: %v", err)
	}
	defer repo.Pool.Close()

	// 3. Inicializar Storage (Supabase)
	storageSvc := storage.NewSupabaseStorage(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		os.Getenv("SUPABASE_BUCKET"),
	)

	// 4. Verificar conexión DB
	if err := repo.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ La DB no responde: %v", err)
	}

	// 5. Inicializar Service (LÓGICA DE NEGOCIO)
	fileService := service.NewFileService(repo, storageSvc)

	// 6. Inicializar Handler (TU CAPA)
	fileHandler := handlers.NewFileHandler(fileService)

	// 7. Configurar Gin
	router := gin.Default()

	// Limitar tamaño de archivos (ej: 10MB)
	router.MaxMultipartMemory = 10 << 20

	// 8. Definir rutas
	router.POST("/upload", fileHandler.Upload)
	router.GET("/download/:token", fileHandler.Download)
	router.GET("/file/:token", fileHandler.GetFile)

	// 9. Levantar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Error al iniciar servidor: %v", err)
	}
}
