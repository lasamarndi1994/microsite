package handler

import (
	"bytes"
	"encoding/json"
	"io"
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
		UserId:       req.UserId,
	}

	var user model.User
	if err := database.DB.First(&user, req.UserId).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("User not found"))
		return
	}

	if err := database.DB.Create(&lead).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to create lead"))
		return
	}

	// Send data to DigiWeb
	partnerCode := user.UserCode
	if partnerCode == "" {
		partnerCode = "ANTIQ"
	}

	go func(l model.Lead, partnerCode string) {
		apiUrl := "https://digiweb.aliceblueonline.com/DigiWebsite"
		payload := map[string]string{
			"leadEmail":  l.Email,
			"leadmobile": l.MobileNumber,
			"Apid":       partnerCode,
			"LeadName":   l.Name,
			"LeadState":  "",
		}

		jsonData, _ := json.Marshal(payload)

		var logEntry model.LeadExternalLog
		logEntry.LeadId = l.Id
		logEntry.ApiUrl = apiUrl
		logEntry.RequestPayload = string(jsonData)

		resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			logEntry.Status = "Failed"
			logEntry.ErrorMessage = err.Error()
			logEntry.HttpStatusCode = 0
		} else {
			defer resp.Body.Close()
			bodyBytes, _ := io.ReadAll(resp.Body)
			logEntry.ResponsePayload = string(bodyBytes)
			logEntry.HttpStatusCode = resp.StatusCode

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				logEntry.Status = "Success"
			} else {
				logEntry.Status = "Failed"
			}
		}

		database.DB.Create(&logEntry)

		// Update lead count in microsite visitor

	}(lead, partnerCode)

	c.JSON(http.StatusCreated, service.SuccessResponse("Thanks for join with me."))
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
