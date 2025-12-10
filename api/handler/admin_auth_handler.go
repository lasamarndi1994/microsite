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

/*
* Handle admin login
* @param c *gin.Context
* @return gin.JSON
 */
func AdminHandleLogin(c *gin.Context) {

	var req request.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Validation errors
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errs,
		})
		return
	}
	var admin model.Admin

	// Find admin by email with active status
	checkAdmin := database.DB.Where("email = ? AND status = ?", req.Email, true).First(&admin)

	// Check if admin exists and password matches
	if checkAdmin.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid email or password"))
		return
	}

	// Use CheckPassword helper for secure password comparison
	if !helper.CheckPassword(admin.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Invalid email or password"))
		return
	}
	token := helper.GenerateToken(10)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login Successfully ",
		"token":   token,
	})
}

/*
* Handle admin registration
* @param c *gin.Context
* @return gin.JSON
 */
func AdminRegister(c *gin.Context) {
	var req request.AdminRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errs,
		})
		return
	}

	var admin model.Admin
	// Check if email already exists
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err == nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Email already exists"))
		return
	}

	// Create new admin
	newAdmin := model.Admin{
		Email:    req.Email,
		Password: helper.HashPassword(req.Password),
		Status:   true,
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to create admin"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Admin registered successfully",
	})
}
