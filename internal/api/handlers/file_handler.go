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
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se encontró el archivo"})
		return
	}
	defer file.Close()

	id := uuid.New()
	ext := filepath.Ext(header.Filename)
	uniqueName := id.String() + ext

	// 1. Subir a Supabase
	path, err := h.storage.Upload(c.Request.Context(), uniqueName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al subir a la nube"})
		return
	}

	// 2. GENERAR EL LINK DE UNA VEZ (Para evitar el GET fallido después)
	signedURL, err := h.storage.GetSignedURL(c.Request.Context(), path)
	if err != nil {
		log.Printf("Error link: %v", err)
	}

	fileMeta := &models.FileMetadata{
		ID:           id,
		Filename:     header.Filename,
		Size:         header.Size,
		ContentType:  header.Header.Get("Content-Type"),
		SupabasePath: path,
		Status:       "pending",
		CreatedAt:    time.Now(),
	}

	// 3. Guardar en Postgres
	err = h.repo.CreateFile(c.Request.Context(), fileMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar en DB"})
		return
	}

	// 4. Devolvemos TODO en el POST
	c.JSON(http.StatusOK, gin.H{
		"message": "Éxito",
		"data": gin.H{
			"id":           id,
			"download_url": signedURL, // <--- Aquí va el link que necesitas
		},
	})
}

func (h *FileHandler) Download(c *gin.Context) {
	idParam := c.Param("id")
	fileID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	fileMeta, err := h.repo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No encontrado"})
		return
	}

	signedURL, err := h.storage.GetSignedURL(c.Request.Context(), fileMeta.SupabasePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error link"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, signedURL)
}
