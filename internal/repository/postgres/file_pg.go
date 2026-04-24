package postgres

import (
	"context"
	"log"
	"os"

	"github.com/JuanSposada/me-transfer/internal/models"

	"github.com/google/uuid"
)

// CreateFile guarda la metadata que le pase la Persona B o C
func (r *PostgresRepo) CreateFile(ctx context.Context, file *models.FileMetadata) error {
	query := `
        INSERT INTO files (id, filename, size, content_type, supabase_path, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	// Ahora enviamos TODOS los campos, incluido el ID que ya generamos en el Handler
	_, err := r.Pool.Exec(ctx, query,
		file.ID,           // $1 (El UUID que ya creamos)
		file.Filename,     // $2
		file.Size,         // $3
		file.ContentType,  // $4
		file.SupabasePath, // $5
		file.Status,       // $6
		file.CreatedAt,    // $7
	)

	return err
}

// GetFileByID busca un archivo por su UUID en la tabla 'files'
func (r *PostgresRepo) GetFileByID(ctx context.Context, id uuid.UUID) (*models.FileMetadata, error) {
	var file models.FileMetadata
	query := `SELECT id, filename, size, content_type, supabase_path, status, created_at 
              FROM files WHERE id = $1`

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
		return nil, err
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

func (r *PostgresRepo) GetExpiredFiles(ctx context.Context) ([]models.FileMetadata, error) {
	var expiredFiles []models.FileMetadata

	query := `
        SELECT id, supabase_path 
        FROM files 
        WHERE created_at < NOW() - INTERVAL '24 hours'`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		log.Printf("DEBUG DB ERROR: %v", err)
		log.Printf("DEBUG URI EN USO: %s", os.Getenv("POSTGRES_URL"))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f models.FileMetadata
		if err := rows.Scan(&f.ID, &f.SupabasePath); err != nil {
			continue
		}
		expiredFiles = append(expiredFiles, f)
	}
	return expiredFiles, nil
}

func (r *PostgresRepo) DeleteFileRecord(ctx context.Context, id uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM files WHERE id = $1", id)
	return err
}
