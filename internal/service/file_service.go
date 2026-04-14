package service

import (
	"context"
	"io"
)

// Interfaces que ya tienes
type StorageService interface {
	Upload(ctx context.Context, filename string, file io.Reader) (string, error)
	GetSignedURL(ctx context.Context, path string) (string, error)
}

type FileRepository interface {
	// Define tus métodos aquí (ej: SaveFile, GetByToken)
}

// Service
type FileService struct {
	repo    FileRepository
	storage StorageService
}

// Constructor
func NewFileService(repo FileRepository, storage StorageService) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

// Lógica de negocio (SIN GIN)
func (s *FileService) UploadFile(ctx context.Context, filename string, file io.Reader) (map[string]interface{}, error) {

	path, err := s.storage.Upload(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	// Aquí luego guardarías en DB (repo)

	return map[string]interface{}{
		"path": path,
	}, nil
}

func (s *FileService) GetDownloadURL(ctx context.Context, token string) (string, error) {

	// Simulación: luego lo sacas de DB
	path := token

	return s.storage.GetSignedURL(ctx, path)
}

func (s *FileService) GetFileByToken(ctx context.Context, token string) (map[string]interface{}, error) {

	// Simulación: luego lo sacas de DB
	return map[string]interface{}{
		"token": token,
	}, nil
}
