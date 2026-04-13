package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/JuanSposada/me-transfer/internal/models"
	"github.com/JuanSposada/me-transfer/internal/repository/postgres"
	"github.com/JuanSposada/me-transfer/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	repo    *postgres.PostgresRepo
	storage *storage.SupabaseStorage
}

func NewFileHandler(repo *postgres.PostgresRepo, storage *storage.SupabaseStorage) *FileHandler {
	return &FileHandler{
		repo:    repo,
		storage: storage,
	}
}

func (h *FileHandler) Upload(c *gin.Context) {
	// 1. Obtener el archivo del formulario
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el archivo en el formulario"})
		return
	}
	defer file.Close()

	// 2. Generar IDs y nombres únicos
	id := uuid.New()
	ext := filepath.Ext(header.Filename)
	uniqueName := id.String() + ext

	// 3. Subir a Supabase (Persona B)
	path, err := h.storage.Upload(c.Request.Context(), uniqueName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir a la nube"})
		return
	}

	// 4. Mapear al modelo exacto que definiste en models/file.go
	fileMeta := &models.FileMetadata{
		ID:           id,
		Filename:     header.Filename, // Nombre original
		Size:         header.Size,     // Tamaño del archivo en bytes
		ContentType:  header.Header.Get("Content-Type"),
		SupabasePath: path,      // La ruta que devuelve Supabase
		Status:       "pending", // Estado inicial
		CreatedAt:    time.Now(),
	}

	// 5. Guardar el registro en Postgres
	err = h.repo.CreateFile(c.Request.Context(), fileMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar en base de datos"})
		return
	}

	// 6. Respuesta al cliente
	c.JSON(http.StatusOK, gin.H{
		"message": "Archivo subido y registrado con éxito",
		"data": gin.H{
			"id":       id,
			"filename": fileMeta.Filename,
			"status":   fileMeta.Status,
		},
	})
}

func (h *FileHandler) Download(c *gin.Context) {
	// 1. Extraer el ID de los parámetros de la URL
	idParam := c.Param("id")
	log.Printf("DEBUG: Buscando ID recibido en URL: [%s]", idParam)

	// 2. Validar que el string sea un UUID válido
	fileID, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("DEBUG: Error parseando UUID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "El formato del ID de archivo es inválido"})
		return
	}

	// 3. Buscar los metadatos en la base de datos (Postgres)
	fileMeta, err := h.repo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		log.Printf("DEBUG: Error en DB para ID %s: %v", fileID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo no encontrado en nuestros registros"})
		return
	}

	log.Printf("DEBUG: Archivo localizado en DB: %s", fileMeta.Filename)

	// 4. Solicitar la URL firmada a Supabase (Persona B)
	// Usamos fileMeta.SupabasePath que es el nombre/ruta que guardamos al subir
	signedURL, err := h.storage.GetSignedURL(c.Request.Context(), fileMeta.SupabasePath)
	if err != nil {
		log.Printf("DEBUG: Error generando URL firmada: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el enlace de descarga"})
		return
	}

	log.Printf("DEBUG: URL firmada generada con éxito para: %s", fileMeta.Filename)

	// 5. Responder al cliente
	// Puedes elegir entre enviar el JSON con el link o redirigir directamente
	c.JSON(http.StatusOK, gin.H{
		"message":      "Enlace generado con éxito",
		"filename":     fileMeta.Filename,
		"size":         fileMeta.Size,
		"content_type": fileMeta.ContentType,
		"download_url": signedURL,
		"expires_in":   "3600 seconds (1 hour)",
	})

	// TIP: Si quieres que el navegador lo descargue apenas entren al link,
	// podrías usar c.Redirect(http.StatusTemporaryRedirect, signedURL)
}
