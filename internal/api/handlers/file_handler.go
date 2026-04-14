package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/JuanSposada/me-transfer/internal/service"
)

type FileHandler struct {
	service *service.FileService
}

// Constructor (IMPORTANTE)
func NewFileHandler(s *service.FileService) *FileHandler {
	return &FileHandler{
		service: s,
	}
}

// POST /upload
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer openedFile.Close()

	result, err := h.service.UploadFile(
		c.Request.Context(),
		file.Filename,
		openedFile,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GET /download/:token
func (h *FileHandler) Download(c *gin.Context) {
	token := c.Param("token")

	url, err := h.service.GetDownloadURL(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": url,
	})
}

// GET /file/:token
func (h *FileHandler) GetFile(c *gin.Context) {
	token := c.Param("token")

	file, err := h.service.GetFileByToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}
