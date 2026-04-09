package storage

import (
	"context"
	"io"
)

// StorageService es la interfaz para manejar archivos en la nube
type StorageService interface {
	// Upload recibe el contenido del archivo y devuelve la URL/Path de Supabase
	Upload(ctx context.Context, fileName string, content io.Reader) (string, error)
	
	// GetSignedURL genera un link temporal para descargar el archivo
	GetSignedURL(ctx context.Context, remotePath string) (string, error)
}