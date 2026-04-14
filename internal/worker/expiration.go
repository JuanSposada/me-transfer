package worker

import (
	"context"
	"log"
	"time"

	"github.com/JuanSposada/me-transfer/internal/repository" // Ajusta a tu path real
	"github.com/JuanSposada/me-transfer/internal/storage"    // Ajusta a tu path real
)

type CleanupWorker struct {
	repo repository.FileRepository
	// Cambiamos FileStorage por StorageService:
	storage storage.StorageService
}

func NewCleanupWorker(repo repository.FileRepository, storage storage.StorageService) *CleanupWorker {
	return &CleanupWorker{
		repo:    repo,
		storage: storage,
	}
}

func (w *CleanupWorker) Start(ctx context.Context) {
	// Se ejecuta cada hora
	ticker := time.NewTicker(1 * time.Hour)

	log.Println("🚀 Worker de limpieza iniciado...")

	go func() {
		w.runCleanup(ctx)
		for {
			select {
			case <-ticker.C:
				w.runCleanup(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *CleanupWorker) runCleanup(ctx context.Context) {
	log.Println("🧹 Iniciando ronda de limpieza...")

	expired, err := w.repo.GetExpiredFiles(ctx)
	if err != nil {
		log.Printf("❌ Error buscando archivos expirados: %v", err)
		return
	}

	for _, file := range expired {
		// 1. Borrar de Supabase
		log.Printf("🛠️ Intentando borrar de Supabase: %s", file.SupabasePath)
		if err := w.storage.DeleteFile(ctx, file.SupabasePath); err != nil {
			log.Printf("⚠️ No se pudo borrar de Supabase [%s]: %v", file.SupabasePath, err)

		}

		// 2. Borrar de Postgres
		if err := w.repo.DeleteFileRecord(ctx, file.ID); err != nil {
			log.Printf("⚠️ No se pudo borrar de la DB [%s]: %v", file.ID, err)
			continue
		}

		log.Printf("🗑️ Archivo eliminado con éxito: %s", file.ID)
	}
}
