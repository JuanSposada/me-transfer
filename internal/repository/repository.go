package repository

import (
	"context"

	"github.com/JuanSposada/me-transfer/internal/models"
	"github.com/google/uuid"
)

// FileRepository es la interfaz que define las operaciones de datos.
// Persona C (API) usará esta interfaz para no depender directamente de Postgres.
type FileRepository interface {
	// CreateFile guarda la metadata y actualiza el struct con ID y fecha
	CreateFile(ctx context.Context, file *models.FileMetadata) error

	// GetFileByID busca un archivo por su UUID
	GetFileByID(ctx context.Context, id uuid.UUID) (*models.FileMetadata, error)

	// CreateToken genera un token de descarga para un archivo
	CreateToken(ctx context.Context, token *models.Token) error

	GetExpiredFiles(ctx context.Context) ([]models.FileMetadata, error)

	// DeleteFileRecord elimina el registro de la base de datos
	DeleteFileRecord(ctx context.Context, id uuid.UUID) error
}
