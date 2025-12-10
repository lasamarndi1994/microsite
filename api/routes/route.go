package routes

import (
	"micro-site/api/handler"
	"micro-site/api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

/*
* Setup the router
* @return *gin.Engine
 */
func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.Static("/uploads", "./uploads")

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"I am fine, Server also working fine": 200})
	})
	api.POST("/validate-mobile-number", handler.ValidateMobileNumber)
	api.POST("/auth/login", handler.HandleLogin)

	api.POST("/auth/resend-otp", middleware.RateLimitMiddleware(rate.Every(1*time.Minute), 1), handler.ResentOtp)

	api.Use(middleware.AuthMiddleware())
	api.POST("/update-profile", handler.UpdateProfile)
	api.POST("/update-profile-image", handler.UploadprofileImage)
	api.GET("/microsite/lists", handler.GetMicroSite)
	api.GET("/microsite/details/:id", handler.GetMicrositeDetails)
	api.POST("/microsite/create", handler.CreateMicrosite)
	api.PUT("/microsite/update/:id", handler.UpdateMicrosite)
	api.DELETE("/microsite/delete/:id", handler.DeleteMicrosite)

	api.POST("/lead/create", handler.CreateLead)
	api.GET("/microsite/:id/leads", handler.GetLeads)
	api.GET("/microsite/view/:slug1/:slug2", handler.GetMicrositeSlugDetails)

	admin := router.Group("/api/admin")
	admin.POST("/login", handler.AdminHandleLogin)

	// Protected admin routes
	admin.Use(middleware.AdminAuthMiddleware())
	admin.GET("/users", handler.GetAllUsers)
	admin.GET("/users/:id/microsites", handler.GetUserMicrosites)
	admin.GET("/microsites", handler.GetAllMicrosites)
	admin.PUT("/microsite/approve/:id", handler.ApproveMicrosite)
	admin.PUT("/microsite/reject/:id", handler.RejectMicrosite)
	admin.DELETE("/microsite/delete/:id", handler.AdminDeleteMicrosite)

	// api.POST("/update-profile", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "Authorized user!"})
	// })

	return router

}
