package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // O el router que prefieran
)

func main() {
	router := gin.Default()

	// Configurar límite de memoria para uploads (ej: 32 MiB)
	router.MaxMultipartMemory = 32 << 20 

	// Grupo de rutas API
	api := router.Group("/api")
	{
		api.POST("/upload", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Upload endpoint ready"})
		})
		api.GET("/download/:token", func(c *gin.Context) {
			token := c.Param("token")
			c.JSON(http.StatusOK, gin.H{"token": token})
		})
		api.GET("/file/:token", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "info"})
		})
	}

	log.Println("Servidor iniciado en :8080")
	router.Run(":8080")
}