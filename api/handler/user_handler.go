package handler

import (
	"micro-site/api/model"
	"micro-site/api/request"
	"micro-site/database"
	"micro-site/internal/helper"
	"micro-site/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
* Update user profile
* @param c *gin.Context
* @return gin.JSON
 */
func UpdateProfile(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	request := request.UserRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}
	user.UserName = request.UserName
	user.AboutMe = request.AboutMe
	user.BusinessName = request.BusinessName
	user.BusinessLocation = request.BusinessLocation
	if request.UserAvatar != "" {
		file_name := strconv.FormatUint(uint64(user.Id), 10) + request.UserName + ".png"
		if service.UploadBase64Image(request.UserAvatar, file_name, "user") {
			user.UserAvatar = file_name
		} else {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
			return
		}
	}
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Internal server error."))
		return
	}
	c.JSON(http.StatusOK, service.SuccessResponse("Successfully updated."))
}

/*
* Upload user profile image
* @param c *gin.Context
* @return gin.JSON
 */
func UploadprofileImage(c *gin.Context) {

	user := c.MustGet("user").(model.User)
	request := request.UserAvatarRequest{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}
	file_name := strconv.FormatUint(uint64(user.Id), 10) + user.UserName + ".png"
	if service.UploadBase64Image(request.UserAvatar, file_name, "user") {
		user.UserAvatar = file_name
	} else {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
		return
	}
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Internal server error."))
		return
	}
	c.JSON(http.StatusAccepted, service.SuccessResponse("Image upload successfully"))

}

func GetUserAnalytics(c *gin.Context) {
	user := c.MustGet("user").(model.User)

	type UserAnalytics struct {
		UserId         uint64 `json:"user_id"`
		UserName       string `json:"user_name"`
		MicrositeCount int64  `json:"microsite_count"`
		TotalViews     int64  `json:"total_views"`
		TotalLeads     int64  `json:"total_leads"`
	}

	var analytics UserAnalytics
	// Raw SQL query for efficiency
	query := `
		SELECT 
			u.id as user_id, 
			u.user_name,
			COUNT(DISTINCT m.id) as microsite_count,
			COALESCE(SUM(m.view_count), 0) as total_views,
			COUNT(DISTINCT l.id) as total_leads,
			COALESCE(SUM(m.engagement_count), 0) as total_engagement
		FROM users u
		LEFT JOIN micro_sites m ON u.id = m.user_id
		LEFT JOIN leads l ON u.id = l.user_id
		GROUP BY u.id
	`
	if err := database.DB.Raw(query, user.Id).Scan(&analytics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch analytics"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Analytics fetched successfully",
		"data":    analytics,
	})
}
