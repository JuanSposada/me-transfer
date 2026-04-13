package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxAllowedSize limita el tamaño de las peticiones entrantes
func MaxAllowedSize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificamos el tamaño antes de que el cuerpo sea procesado
		if c.Request.ContentLength > limit {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "El archivo es demasiado grande",
				"limit": limit,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
