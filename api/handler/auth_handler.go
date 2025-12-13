package handler

import (
	"micro-site/api/model"
	"micro-site/api/request"
	"micro-site/database"
	"micro-site/internal/helper"
	"micro-site/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/*
* Handle user login
* @param c *gin.Context
* @return gin.JSON
 */
func HandleLogin(c *gin.Context) {
	input := request.LoginRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		// Validation errors
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errs,
		})
		return
	}
	user := model.User{}
	mobileCheck := database.DB.Where("mobile_number =? AND status = ? ", input.MobileNumber, "Active").First(&user)
	if mobileCheck.RowsAffected == 0 {
		c.JSON(http.StatusOK, service.ErrorResponse("Enter mobile number is invalid."))
		return
	}

	// Check if password is provided
	if input.Password != "" {
		if !helper.CheckPassword(user.Password, input.Password) {
			c.JSON(http.StatusBadRequest, service.ErrorResponse("Enter credentials are invalid."))
			return
		}
		// JWT token genreate
		token, err := service.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("Unable to generate token"))
			return
		}
		c.JSON(http.StatusOK, service.SuccessResponse("Login Successfully", token))
		return
	}

	// Fallback to OTP check
	if input.MobileOtp == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("OTP is required if password is not provided"))
		return
	}

	otpModel := model.Otp{}
	otpCheck := database.DB.Where("user_id =? ", user.Id).Where("email_otp = ? OR whatapp_otp = ?", input.MobileOtp, input.MobileOtp).Order("created_at DESC").First(&otpModel)

	if otpCheck.RowsAffected > 0 {
		expireTime := otpModel.CreatedAt.Add(5 * time.Minute)
		if time.Now().After(expireTime) || otpModel.Status {
			c.JSON(http.StatusBadRequest, service.ErrorResponse("Enter OTP is expired."))
			return
		}
		otpModel.Status = true
		database.DB.Save(&otpModel)
		// JWT token genreate
		token, err := service.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("Unable to generate token"))
			return
		}
		c.JSON(http.StatusOK, service.SuccessResponse("Login Successfully", token))

		return
	} else {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Enter OTP is Invalid"))
		return
	}
}

/*
* Update user password
* @param c *gin.Context
* @return gin.JSON
 */
func UpdatePassword(c *gin.Context) {
	var req request.UpdatePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}

	user := model.User{}
	mobileCheck := database.DB.Where("mobile_number =? ", req.MobileNumber).First(&user)
	if mobileCheck.RowsAffected == 0 {
		c.JSON(http.StatusOK, service.ErrorResponse("Enter mobile number is invalid."))
		return
	}

	// Hash the new password
	hashedPassword := helper.HashPassword(req.Password)
	user.Password = hashedPassword

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to update password"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse("Your account created successfully.Please relogin process"))
}

/*
* Validate mobile number
* @param c *gin.Context
* @return gin.JSON
 */
func ValidateMobileNumber(c *gin.Context) {
	request := request.MobileRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}
	user := model.User{}
	result := database.DB.Where("mobile_number = ?", request.MobileNumber).First(&user)

	if result.RowsAffected > 0 {
		if user.Password != "" {
			c.JSON(http.StatusOK, service.SuccessResponse("Password is already set.", true))
			return
		}

		c.JSON(http.StatusOK, service.SuccessResponse("Password is not set.", false))

		// send otp
		SendOtp(user)
		user.Email = helper.MaskEmail(user.Email)
		// c.JSON(http.StatusOK, service.SuccessResponse("OTP is send your email address.", user))
		return
	} else {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Enter mobile number doesn't exist."))
		return
	}
}

/*
* Resend OTP to user
* @param c *gin.Context
* @return gin.JSON
 */
func ResentOtp(c *gin.Context) {
	request := request.MobileRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}
	user := model.User{}
	result := database.DB.Where("mobile_number = ?", request.MobileNumber).First(&user)

	if result.RowsAffected > 0 {
		SendOtp(user)
		c.JSON(http.StatusOK, service.SuccessResponse("OTP is resend your register email address.", helper.MaskEmail(user.Email)))
		return
	}
}

/*
* Send OTP to user
* @param user model.User
* @return void
 */
func SendOtp(user model.User) {

	email_otp, _ := helper.GenerateOTP()
	whatapp_otp, _ := helper.GenerateOTP()

	//save otp
	otpModel := model.Otp{}
	otpModel.EmailOtp = int32(email_otp)
	otpModel.WhatappOtp = int32(whatapp_otp)
	otpModel.UserId = user.Id
	database.DB.Save(&otpModel)

	// go service.SendHTMLEmail(user.Email, "Welcome !", service.EmailData{
	// 	Name:      user.UserName,
	// 	Email:     user.Email,
	// 	OtpNumber: int32(email_otp),
	// })
}

func GetAuthUserDetails(c *gin.Context) {
	user := c.MustGet("user").(model.User)

	var count int64
	database.DB.Model(&model.MicroSite{}).Where("user_id = ?", user.Id).Count(&count)
	user.MicrositeCount = count

	c.JSON(http.StatusOK, service.SuccessResponse("User details fetched successfully", user))
}
