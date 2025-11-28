package handler

import (
	"micro-site/api/model"
	"micro-site/api/request"
	"micro-site/database"
	"micro-site/internal/helper"
	"micro-site/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateMobileNumber(c *gin.Context) {
	type PhoneNumber struct {
		MobileNumber string `json:"mobile_number" binding:"required"`
	}

	var mobile_number PhoneNumber
	if err := c.ShouldBindJSON(&mobile_number); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
	}
	var user model.User
	result := database.DB.Where("mobile_number = ?", mobile_number.MobileNumber).First(&user)
	if result.RowsAffected > 0 {
		// send otp
		err := service.SendSMS(service.NewSNSClient(), "+91"+mobile_number.MobileNumber, "Your OTP is 1234")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "OTP sent"})

		return
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": "Mobile number doest exists",
		})
		return
	}
}
func HandleLogin(c *gin.Context) {
	var input request.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		// Validation errors
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"errors": errs,
		})
		return
	}
}
