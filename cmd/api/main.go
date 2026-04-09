package main

import (
	"context"
	"log"
	"os"

	"github.com/JuanSposada/me-transfer/internal/repository/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se encontró archivo .env, usando variables de sistema")
	}

	// 2. Obtener URL de conexión
	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		log.Fatal("❌ La variable POSTGRES_URL no está definida")
	}

	// 3. Inicializar el Repositorio de Postgres (Tu trabajo previo)
	repo, err := postgres.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la DB: %v", err)
	}
	defer repo.Pool.Close()

	// 4. Verificar conexión (Ping)
	if err := repo.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ La DB no responde: %v", err)
	}

	log.Println("✅ Infraestructura lista: PostgreSQL conectado correctamente")

	// --- AQUÍ ENTRARÁN LA PERSONA B Y C ---
	// Persona B: storageService := supabase.NewStorage(...)
	// Persona C: server := api.NewServer(repo, storageService)
	// server.Start()
}
