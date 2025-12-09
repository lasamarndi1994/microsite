package handler

import (
	"fmt"
	"micro-site/api/model"
	"micro-site/api/request"
	"micro-site/database"
	"micro-site/internal/helper"
	"micro-site/internal/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
* Get all microsites for the authenticated user
* @param c *gin.Context
* @return gin.JSON
 */
func GetMicroSite(c *gin.Context) {
	// Get the authenticated user
	data, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("User not authenticated"))
		return
	}
	user, ok := data.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Invalid user data"))
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	// Build base query
	query := database.DB.Model(&model.MicroSite{}).Where("user_id = ?", user.Id)

	// Get status filter from query parameter
	status := c.Query("status")

	// Apply status filter if provided
	if status != "" {
		// Validate status value
		validStatuses := map[string]bool{
			"Pending":  true,
			"Approved": true,
			"Rejected": true,
			"Active":   true,
		}
		if !validStatuses[status] {
			c.JSON(http.StatusBadRequest, service.ErrorResponse("Invalid status. Valid values are: Pending, Approved, Rejected, Active"))
			return
		}
		query = query.Where("status = ?", status)
	}

	// Count total records
	var total int64
	query.Count(&total)

	// Fetch paginated data
	var microsites []model.MicroSite
	if err := query.Preload("Services").
		Preload("SocialLinks").
		Limit(limit).
		Offset(offset).
		Find(&microsites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch microsites"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "Microsites fetched successfully",
		"data":        microsites,
		"total":       total,
		"page":        page,
		"per_page":    limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

/*
* Get details of a specific microsite
* @param c *gin.Context
* @return gin.JSON
 */
func GetMicrositeDetails(c *gin.Context) {
	// Get the authenticated user
	data, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("User not authenticated"))
		return
	}
	user, ok := data.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Invalid user data"))
		return
	}
	// Get microsite ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	// Fetch microsite with related data
	var microsite model.MicroSite
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.Id).
		Preload("Services").
		Preload("SocialLinks").
		First(&microsite).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("Microsite not found or you don't have permission to access it"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Microsite details fetched successfully",
		"data":    microsite,
	})
}

/*
* Create a new microsite
* @param c *gin.Context
* @return gin.JSON
 */
func CreateMicrosite(c *gin.Context) {
	data, _ := c.Get("user")
	user, _ := data.(model.User)
	var req request.MicrositeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Validation errors
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}
	var micro_site model.MicroSite

	micro_site.UserId = user.Id
	micro_site.Title = req.Title
	micro_site.FullName = req.FullName
	micro_site.Description = req.Description
	micro_site.IsDraft = req.IsDraft

	randName := fmt.Sprintf("%d", os.Getpid())
	if req.AvatarIcon != "" {
		file_name := strconv.FormatUint(uint64(user.Id), 10) + user.UserName + randName + ".png"
		if service.UploadBase64Image(req.AvatarIcon, file_name, "avatar") {
			micro_site.AvatarIcon = file_name
		} else {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
			return
		}
	}
	if req.BannerImage != "" {
		file_name := strconv.FormatUint(uint64(user.Id), 10) + user.UserName + randName + ".png"
		if service.UploadBase64Image(req.BannerImage, file_name, "banner") {
			micro_site.BannerImage = file_name
		} else {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
			return
		}
	}
	// Convert service
	for _, serviceName := range req.ServicesName {
		micro_site.Services = append(micro_site.Services, model.Service{
			UserId: user.Id,
			Name:   serviceName.Name,
		})
	}

	// Convert social links
	for _, social := range req.SocialLink {
		micro_site.SocialLinks = append(micro_site.SocialLinks, model.SocialLink{
			UserId: user.Id,
			Type:   social.Type,
			Url:    social.Url,
		})
	}

	if err := database.DB.Create(&micro_site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Internal server error."))
		return
	}
	c.JSON(http.StatusAccepted, service.SuccessResponse("Successfully created your micro site."))
}

/*
* Update an existing microsite
* @param c *gin.Context
* @return gin.JSON
 */
func UpdateMicrosite(c *gin.Context) {
	data, _ := c.Get("user")
	user, _ := data.(model.User)
	id := c.Param("id")

	var existing model.MicroSite
	if err := database.DB.Preload("Services").Preload("SocialLinks").First(&existing, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Microsite not found"})
		return
	}

	var req request.MicrositeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := helper.FormatValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"errors": errs,
		})
		return
	}

	// Update Parent Data
	existing.FullName = req.FullName
	existing.Title = req.Title
	existing.Description = req.Description
	existing.IsDraft = req.IsDraft

	randName := fmt.Sprintf("%d", os.Getpid())
	if req.AvatarIcon != "" {
		file_name := strconv.FormatUint(uint64(user.Id), 10) + user.UserName + randName + ".png"
		if service.UploadBase64Image(req.AvatarIcon, file_name, "avatar") {
			existing.AvatarIcon = file_name
		} else {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
			return
		}
	}
	if req.BannerImage != "" {
		file_name := strconv.FormatUint(uint64(user.Id), 10) + user.UserName + randName + ".png"
		if service.UploadBase64Image(req.BannerImage, file_name, "banner") {
			existing.BannerImage = file_name
		} else {
			c.JSON(http.StatusInternalServerError, service.ErrorResponse("File upload failed"))
			return
		}
	}

	// DELETE Children First

	database.DB.Where("micro_site_id = ?", id).Delete(&model.Service{})
	database.DB.Where("micro_site_id = ?", id).Delete(&model.SocialLink{})
	existing.Services = []model.Service{}
	existing.SocialLinks = []model.SocialLink{}

	// ADD Updated Service list
	for _, serviceName := range req.ServicesName {
		existing.Services = append(existing.Services, model.Service{
			Name:        serviceName.Name,
			MicroSiteId: existing.Id,
		})
	}

	// Convert social links
	for _, social := range req.SocialLink {
		existing.SocialLinks = append(existing.SocialLinks, model.SocialLink{
			MicroSiteId: existing.Id,
			Type:        social.Type,
			Url:         social.Url,
		})
	}
	database.DB.Save(&existing)
	c.JSON(http.StatusOK, service.SuccessResponse("Microsite updated successfully"))
}

/*
* Delete a microsite
* @param c *gin.Context
* @return gin.JSON
 */
func DeleteMicrosite(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&model.MicroSite{}, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Failed to delete microsite"))
		return
	}
	c.JSON(http.StatusAccepted, service.SuccessResponse("Microsite deleted successfully"))
}
