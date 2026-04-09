package storage

import (
	"context"
	"io"
	"fmt"
)

type SupabaseStorage struct {
	URL    string
	Key    string
	Bucket string
}

func NewSupabaseStorage(url, key, bucket string) *SupabaseStorage {
	return &SupabaseStorage{
		URL:    url,
		Key:    key,
		Bucket: bucket,
	}
}

func (s *SupabaseStorage) Upload(ctx context.Context, fileName string, content io.Reader) (string, error) {
	// TODO: Persona B implementará la llamada a la API de Supabase aquí
	return fmt.Sprintf("paths/%s", fileName), nil
}

func (s *SupabaseStorage) GetSignedURL(ctx context.Context, remotePath string) (string, error) {
	// TODO: Persona B implementará la generación de URLs firmadas aquí
	return "https://supabase.com/fake-signed-url", nil
}