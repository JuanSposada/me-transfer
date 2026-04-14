package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SupabaseStorage struct {
	url    string
	key    string
	bucket string
}

func NewSupabaseStorage(url, key, bucket string) *SupabaseStorage {
	return &SupabaseStorage{
		url:    strings.TrimRight(url, "/"),
		key:    key,
		bucket: bucket,
	}
}

func (s *SupabaseStorage) Upload(ctx context.Context, fileName string, content io.Reader) (string, error) {
	// Construimos la URL manualmente: /storage/v1/object/nombre-bucket/archivo
	fullURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, s.bucket, fileName)

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, content)
	if err != nil {
		return "", fmt.Errorf("error creando petición: %v", err)
	}

	// Headers obligatorios para Supabase
	req.Header.Set("Authorization", "Bearer "+s.key)
	req.Header.Set("apiKey", s.key)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error enviando archivo (red): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase respondió con código %d: %s", resp.StatusCode, string(body))
	}

	return fileName, nil
}

func (s *SupabaseStorage) GetSignedURL(ctx context.Context, fileName string) (string, error) {
	fullURL := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", s.url, s.bucket, fileName)
	body := strings.NewReader(`{"expiresIn": 3600}`)

	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.key)
	req.Header.Set("apiKey", s.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// --- 1. Definimos la estructura para capturar el JSON ---
	var result struct {
		SignedURL string `json:"signedURL"`
	}

	// --- 2. Decodificamos el cuerpo de la respuesta en la variable 'result' ---
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decodificando respuesta de Supabase: %v", err)
	}

	// 1. Limpiamos la URL base (quitamos barras al final)
	baseURL := strings.TrimRight(s.url, "/")

	// 2. Verificamos si result.SignedURL ya trae el prefijo /storage/v1
	// Supabase a veces devuelve solo /object/sign/...
	signedPath := result.SignedURL
	if !strings.HasPrefix(signedPath, "/storage/v1") {
		signedPath = "/storage/v1" + signedPath
	}

	// 3. Unimos todo
	fullSignedURL := baseURL + signedPath

	return fullSignedURL, nil
}

func (s *SupabaseStorage) DeleteFile(ctx context.Context, supabasePath string) error {
	// La URL para borrado múltiple es correcta
	fullURL := fmt.Sprintf("%s/storage/v1/object/%s", s.url, s.bucket)

	// CAMBIO AQUÍ: Usamos "prefixes" pero asegúrate de que el path sea el exacto
	// Si sigue fallando, la API de Supabase prefiere "filenames" para archivos individuales
	payload := fmt.Sprintf(`{"prefixes": ["%s"]}`, supabasePath)
	body := strings.NewReader(payload)

	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.key)
	req.Header.Set("apiKey", s.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// IMPORTANTE: Supabase devuelve un array de objetos en el body indicando
	// qué se borró y qué no. Vamos a leerlo para debuguear si falla.
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error al borrar en supabase (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
