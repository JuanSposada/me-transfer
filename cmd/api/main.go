package main

import (
	"context"
	"log"
	"os"

	"github.com/JuanSposada/me-transfer/internal/repository/postgres"
	"github.com/JuanSposada/me-transfer/internal/storage" // Asegúrate de que este import exista
	"github.com/joho/godotenv"
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
	defer repo.Pool.Close() // Importante para cerrar conexiones al apagar

	// 3. Inicializar Storage (Supabase)
	// Aquí inyectamos las llaves que vas a crear ahora
	storageSvc := storage.NewSupabaseStorage(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		os.Getenv("SUPABASE_BUCKET"),
	)
	_ = storageSvc // Para evitar error de variable no usada, lo usaremos en Persona B

	// 4. Verificación final
	if err := repo.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ La DB no responde: %v", err)
	}

	log.Println("✅ INFRAESTRUCTURA COMPLETA: Postgres y Supabase configurados.")

	// Bloqueo para que no se cierre el programa (Temporal hasta que Persona C ponga la API)
	select {}
}
