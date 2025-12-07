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
* Handle Login
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
	mobileCheck := database.DB.Where("mobile_number =? ", input.MobileNumber).First(&user)
	if mobileCheck.RowsAffected == 0 {
		c.JSON(http.StatusOK, service.ErrorResponse("Enter mobile number is invalid."))
		return
	}
	otpModel := model.Otp{}
	otpCheck := database.DB.Where("user_id =? ", user.Id).Where("email_otp = ? OR whatapp_otp = ?", input.MobileOtp, input.MobileOtp).Order("created_at DESC").First(&otpModel)

	if otpCheck.RowsAffected > 0 {
		expireTime := otpModel.CreatedAt.Add(5 * time.Minute)
		if time.Now().After(expireTime) || otpModel.Status {
			c.JSON(http.StatusOK, service.ErrorResponse("Enter OTP is expired."))
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
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Login Successfully ",
			"token":   token,
		})
		return
	} else {
		c.JSON(http.StatusOK, service.ErrorResponse("Enter OTP is Invalid"))
		return
	}
}

/*
* validate mobile number
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
		// send otp
		SendOtp(user)
		c.JSON(http.StatusOK, service.SuccessResponse("OTP is send your email address.", helper.MaskEmail(user.Email)))
		return
	} else {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Enter mobile number doesn't exist."))
		return
	}
}

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

func SendOtp(user model.User) {

	email_otp, _ := helper.GenerateOTP()
	whatapp_otp, _ := helper.GenerateOTP()

	//save otp
	otpModel := model.Otp{}
	otpModel.EmailOtp = int32(email_otp)
	otpModel.WhatappOtp = int32(whatapp_otp)
	otpModel.UserId = user.Id
	database.DB.Save(&otpModel)

	go service.SendHTMLEmail(user.Email, "Welcome !", service.EmailData{
		Name:      user.UserName,
		Email:     user.Email,
		OtpNumber: int32(email_otp),
	})
}
