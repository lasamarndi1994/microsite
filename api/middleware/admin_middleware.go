package middleware

import (
	"micro-site/api/model"
	"micro-site/config"
	"micro-site/database"
	"micro-site/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*
* Middleware to validate admin authentication
* @return gin.HandlerFunc
 */
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		tokenString := authHeader
		if len(tokenString) > 7 && strings.ToUpper(tokenString[:7]) == "BEARER " {
			tokenString = tokenString[7:]
		}

		secret := config.LoadConfig().JWTSecretKey
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid token claims"))
			c.Abort()
			return
		}

		// Check if user_id exists in claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid token data"))
			c.Abort()
			return
		}
		adminID := uint64(userIDFloat)

		// Validate admin token (check if admin exists and is active)
		var admin model.Admin
		if err := database.DB.Where("id = ? AND status = ?", adminID, true).First(&admin).Error; err != nil {
			c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid admin credentials"))
			c.Abort()
			return
		}

		// Store admin in context for use in handlers
		c.Set("admin", admin)
		c.Next()
	}
}
