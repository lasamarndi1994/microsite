package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BodySizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set your max size (e.g., 50 MB)
		maxBytes := int64(50 * 1024 * 1024) // 50MB
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
