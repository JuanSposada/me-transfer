package storage

import (
	"context"
	"io"
)

// StorageService es la interfaz para manejar archivos en la nube
type StorageService interface {
	Upload(ctx context.Context, fileName string, content io.Reader) (string, error)
	GetSignedURL(ctx context.Context, remotePath string) (string, error)
	// Añadimos esto para que el Worker pueda limpiar:
	DeleteFile(ctx context.Context, remotePath string) error
}
