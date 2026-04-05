package postgres

import (
	"context"

	"github.com/JuanSposada/me-transfer/internal/models"
)

// CreateFile guarda la metadata que le pase la Persona B o C
func (r *PostgresRepo) CreateFile(ctx context.Context, file *models.FileMetadata) error {
	query := `
		INSERT INTO files (filename, size, content_type, supabase_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, status
	`

	// Ejecutamos y "escaneamos" los valores que genera la DB (ID y Fecha)
	// de vuelta a nuestro struct de Go
	return r.Pool.QueryRow(ctx, query,
		file.Filename,
		file.Size,
		file.ContentType,
		file.SupabasePath,
	).Scan(&file.ID, &file.CreatedAt, &file.Status)
}
