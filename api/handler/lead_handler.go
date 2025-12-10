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
* Create a new lead
* @param c *gin.Context
* @return gin.JSON
 */
func CreateLead(c *gin.Context) {
	var req request.LeadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}

	lead := model.Lead{
		Name:         req.Name,
		Email:        req.Email,
		MobileNumber: req.MobileNumber,
		MicrositeId:  req.MicrositeId,
	}

	if err := database.DB.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to create lead"))
		return
	}

	c.JSON(http.StatusCreated, service.SuccessResponse("Lead created successfully", lead))
}

/*
* Get leads for a specific microsite
* @param c *gin.Context
* @return gin.JSON
 */
func GetLeads(c *gin.Context) {
	micrositeId := c.Param("id")
	if micrositeId == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	var leads []model.Lead
	if err := database.DB.Where("microsite_id = ?", micrositeId).Find(&leads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch leads"))
		return
	}

	c.JSON(http.StatusOK, service.SuccessResponse("Thanks for join with me", leads))
}
