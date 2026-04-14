package api

import (
    "github.com/gin-gonic/gin"
    "github.com/JuanSposada/me-transfer/internal/api/handlers"
)

func NewRouter(handler *handlers.FileHandler) *gin.Engine {
    r := gin.Default()

    // Rutas
    r.POST("/upload", handler.Upload)
    r.GET("/download/:token", handler.Download)
    r.GET("/file/:token", handler.GetFile)

    return r
}