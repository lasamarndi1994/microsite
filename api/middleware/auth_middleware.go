package middleware

import (
	"micro-site/api/model"
	"micro-site/config"
	"micro-site/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*
* Middleware to authenticate user via JWT
* @return gin.HandlerFunc
 */
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Authorization required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		//cfg := config.LoadConfig()
		secret := config.LoadConfig().JWTSecretKey

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}
		// Extract user_id and store in context
		claims := token.Claims.(jwt.MapClaims)

		userID := uint64(claims["user_id"].(float64))

		user := model.User{}

		//database.DB.Where("id = ?", userID).First(&user)
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}
		if user.Id != userID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
