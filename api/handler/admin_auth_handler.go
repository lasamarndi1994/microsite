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
