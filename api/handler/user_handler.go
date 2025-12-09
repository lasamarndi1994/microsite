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
