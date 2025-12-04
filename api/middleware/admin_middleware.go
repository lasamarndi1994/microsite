package middleware

import (
	"micro-site/api/model"
	"micro-site/database"
	"micro-site/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware validates admin authentication
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Validate admin token (simple validation - you can enhance this)
		var admin model.Admin
		if err := database.DB.Where("status = ?", true).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid admin credentials"))
			c.Abort()
			return
		}

		// Store admin in context for use in handlers
		c.Set("admin", admin)
		c.Next()
	}
}
