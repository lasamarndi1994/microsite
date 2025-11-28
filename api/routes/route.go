package routes

import (
	"micro-site/api/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"I am fine": 200})
	})
	api.POST("/validate-mobile-number", handler.ValidateMobileNumber)
	api.POST("/auth/login", handler.HandleLogin)
	return router

}
