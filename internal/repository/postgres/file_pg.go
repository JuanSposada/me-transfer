package postgres

import (
	"context"

	"github.com/JuanSposada/me-transfer/internal/models"

	"fmt"
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

// GetFileByID busca un archivo por su UUID en la tabla 'files'
func (r *PostgresRepo) GetFileByID(ctx context.Context, id string) (*models.FileMetadata, error) {
	query := `
		SELECT id, filename, size, content_type, supabase_path, status, created_at
		FROM files
		WHERE id = $1
	`

	var file models.FileMetadata
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.Filename,
		&file.Size,
		&file.ContentType,
		&file.SupabasePath,
		&file.Status,
		&file.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("archivo no encontrado: %w", err)
	}

	return &file, nil
}

// CreateToken inserta un nuevo token vinculado a un archivo
func (r *PostgresRepo) CreateToken(ctx context.Context, token *models.Token) error {
	query := `
		INSERT INTO tokens (file_id, expires_at)
		VALUES ($1, $2)
		RETURNING token, created_at
	`

	// Al igual que con el archivo, la DB genera el UUID del token y la fecha
	return r.Pool.QueryRow(ctx, query,
		token.FileID,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}
