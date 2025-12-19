// middleware/activity_logger.go
package middleware

import (
	"micro-site/api/model"
	"micro-site/database"
	"micro-site/internal/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

func ActivityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/uploads" {
			c.Next()
			return
		}
		start := time.Now()

		c.Next() // process request

		go func() {
			duration := time.Since(start)

			ua := user_agent.New(c.Request.UserAgent())
			browserName, _ := ua.Browser()

			activity := model.UserActivity{
				UserID:       helper.GetUserID(c), // nullable
				Method:       c.Request.Method,
				Path:         c.Request.URL.Path,
				IP:           helper.GetClientIP(c),
				UserAgent:    browserName,
				StatusCode:   c.Writer.Status(),
				ResponseTime: duration.Milliseconds(),
			}

			// Save asynchronously (non-blocking)
			database.DB.Create(&activity)
		}()
	}
}
