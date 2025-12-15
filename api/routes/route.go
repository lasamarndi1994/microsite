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
	router.Use(middleware.BodySizeMiddleware())

	// router.Use(middleware.ActivityLogger())

	router.Static("/api/uploads", "./uploads")

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"I am fine, Server also working fine": 200})
	})
	api.POST("/validate-mobile-number", handler.ValidateMobileNumber)
	api.POST("/auth/login", handler.HandleLogin)
	api.POST("/auth/update-password", handler.UpdatePassword)
	api.POST("/auth/resend-otp", middleware.RateLimitMiddleware(rate.Every(1*time.Minute), 1), handler.ResentOtp)
	api.POST("/auth/forgot-password", handler.ForgotPassword)
	api.POST("/auth/reset-password", handler.ResetPassword)
	api.GET("/microsite/view/:slug1/:slug2", handler.GetMicrositeSlugDetails)
	api.POST("/lead/create", handler.CreateLead)
	api.POST("/microsite/engagement/:slug", handler.UpdateMicrositeEngagementCount)

	api.Use(middleware.AuthMiddleware())

	api.GET("/auth/user", handler.GetAuthUserDetails)
	api.POST("/update-profile", handler.UpdateProfile)
	api.POST("/update-profile-image", handler.UploadprofileImage)
	api.GET("/microsite/lists", handler.GetMicroSite)
	api.GET("/microsite/details/:uuid", handler.GetMicrositeDetails)
	api.POST("/microsite/create", handler.CreateMicrosite)
	api.PUT("/microsite/update/:uuid", handler.UpdateMicrosite)
	api.DELETE("/microsite/delete/:uuid", handler.DeleteMicrosite)
	api.GET("/microsite/search", handler.SearchMicrosite)

	api.GET("/microsite/:id/leads", handler.GetLeads)
	api.GET("/analytics", handler.GetUserAnalytics)

	admin := router.Group("/api/admin")
	admin.POST("/login", handler.AdminHandleLogin)
	admin.POST("/register", handler.AdminRegister)

	// Protected admin routes
	admin.Use(middleware.AdminAuthMiddleware())
	admin.GET("/users", handler.GetAllUsers)
	admin.GET("/users/:uuid/microsites", handler.GetUserMicrosites)

	admin.GET("/microsites", handler.GetAllMicrosites)
	admin.GET("/microsite/:uuid", handler.AdminGetMicrositeDetails)
	admin.PUT("/microsite/approve/:uuid", handler.ApproveMicrosite)
	admin.PUT("/microsite/reject/:uuid", handler.RejectMicrosite)
	admin.DELETE("/microsite/delete/:uuid", handler.AdminDeleteMicrosite)

	return router

}
