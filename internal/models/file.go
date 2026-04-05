package models

import (
	"time"

	"github.com/google/uuid"
)

// FileMetadata representa la tabla 'files'
// Usamos tags de json para la API y db para las queries de Postgres
type FileMetadata struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Filename     string    `db:"filename" json:"filename"`
	Size         int64     `db:"size" json:"size"`
	ContentType  string    `db:"content_type" json:"content_type"`
	SupabasePath string    `db:"supabase_path" json:"supabase_path"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// Token representa la tabla 'tokens'
// Este modelo será clave para generar los enlaces de descarga
type Token struct {
	ID        uuid.UUID `db:"token" json:"token"`
	FileID    uuid.UUID `db:"file_id" json:"file_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
