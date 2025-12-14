// middleware/activity_logger.go
package middleware

import (
	"micro-site/api/model"
	"micro-site/database"
	"micro-site/internal/helper"
	"time"

	"github.com/gin-gonic/gin"
)

func ActivityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // process request

		duration := time.Since(start)

		activity := model.UserActivity{
			UserID:       helper.GetUserID(c), // nullable
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			ResponseTime: duration.Milliseconds(),
		}

		// Save asynchronously (non-blocking)
		go database.DB.Create(&activity)
	}
}
