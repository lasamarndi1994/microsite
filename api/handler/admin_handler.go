package handler

import (
	"math"
	"micro-site/api/model"
	"micro-site/database"
	"micro-site/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
* Display all users list
* @param c *gin.Context
* @return gin.JSON
 */
func GetAllUsers(c *gin.Context) {
	// Get admin from context (set by middleware)
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Optional status filter

	// Build query
	query := database.DB.Model(&model.User{})

	// Search filter
	search := c.Query("search")
	if search != "" {
		searchLike := "%" + search + "%"
		query = query.Where("user_name LIKE ? OR email LIKE ? OR mobile_number LIKE ?", searchLike, searchLike, searchLike)
	}

	// Apply status filter if provided

	// Date filter
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if fromDate != "" && toDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", fromDate, toDate)
	} else if fromDate != "" {
		query = query.Where("created_at >= ?", fromDate)
	} else if toDate != "" {
		query = query.Where("created_at <= ?", toDate)
	}

	// Count total records
	var totalRecords int64
	query.Count(&totalRecords)

	// Execute query with pagination
	var users []model.User
	if err := query.Order("updated_at desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch users"))
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Get microsite counts for fetched users
	if len(users) > 0 {
		var userIDs []uint64
		for _, u := range users {
			userIDs = append(userIDs, u.Id)
		}

		type UserCount struct {
			UserId int64
			Count  int64
		}
		var counts []UserCount

		// Aggregate counts
		database.DB.Model(&model.MicroSite{}).
			Select("user_id, count(*) as count").
			Where("user_id IN ?", userIDs).
			Group("user_id").
			Scan(&counts)

		// Map counts to users
		countMap := make(map[uint64]int64)
		for _, c := range counts {
			countMap[uint64(c.UserId)] = c.Count
		}

		for i := range users {
			users[i].MicrositeCount = countMap[users[i].Id]
		}
	}

	// Return success response with pagination
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Users fetched successfully",
		"data":    users,
		"pagination": gin.H{
			"current_page":  page,
			"limit":         limit,
			"total_records": totalRecords,
			"total_pages":   totalPages,
		},
	})
}

/*
* Display all microsites for a specific user
* @param c *gin.Context
* @return gin.JSON
 */
func GetUserMicrosites(c *gin.Context) {
	// Get admin from context
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Get user UUID from URL parameter
	userUUID := c.Param("uuid")
	if userUUID == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("User UUID is required"))
		return
	}

	// Verify user exists and get User ID
	var user model.User
	if err := database.DB.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("User not found"))
		return
	}

	// Optional status filter for microsites
	status := c.Query("status")

	// Build query for microsites using User ID
	query := database.DB.Where("user_id = ? AND status != ?", user.Id, "Draft").
		Preload("Services").
		Preload("SocialLinks")

	// Apply status filter if provided
	if status != "" {
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

	// Search filter
	search := c.Query("search")
	if search != "" {
		searchLike := "%" + search + "%"
		query = query.Where("title LIKE ? OR sub_title LIKE ?", searchLike, searchLike)
	}

	// Fetch microsites
	var microsites []model.MicroSite
	if err := query.Order("updated_at desc").Find(&microsites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch microsites"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User microsites fetched successfully",
		"user": gin.H{
			"id":            user.Id,
			"uuid":          user.Uuid,
			"user_name":     user.UserName,
			"email":         user.Email,
			"mobile_number": user.MobileNumber,
			"status":        user.Status,
			"user_avatar":   user.UserAvatar,
		},
		"data":  microsites,
		"count": len(microsites),
	})
}

/*
* Display all microsites across all users
* @param c *gin.Context
* @return gin.JSON
 */
func GetAllMicrosites(c *gin.Context) {
	// Get admin from context
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Optional status filter
	status := c.Query("status")

	// Build query
	query := database.DB.Preload("Services").Preload("SocialLinks")

	// Apply status filter if provided
	if status != "" {
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

	// Fetch all microsites
	var microsites []model.MicroSite
	if err := query.Find(&microsites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch microsites"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "All microsites fetched successfully",
		"data":    microsites,
		"count":   len(microsites),
	})
}

/*
* Approve a pending microsite
* @param c *gin.Context
* @return gin.JSON
 */
func ApproveMicrosite(c *gin.Context) {
	// Get admin from context
	data, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}
	admin, ok := data.(model.Admin)
	if !ok {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Invalid admin data"))
		return
	}

	// Get microsite ID from URL parameter
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	// Parse request body for optional comment
	var reqBody struct {
		ApproveMessage string `json:"approve_message"`
	}
	c.ShouldBindJSON(&reqBody)

	// Fetch microsite
	var microsite model.MicroSite
	if err := database.DB.Where("uuid = ?", uuid).First(&microsite).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("Microsite not found"))
		return
	}

	// Update microsite status
	microsite.Status = "Approved"
	microsite.AdminId = admin.Id
	if reqBody.ApproveMessage != "" {
		microsite.ApproveMessage = reqBody.ApproveMessage
	}

	if err := database.DB.Save(&microsite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to approve microsite"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Microsite approved successfully",
		"data":    microsite,
	})
}

/*
* Reject a microsite with comment
* @param c *gin.Context
* @return gin.JSON
 */
func RejectMicrosite(c *gin.Context) {
	// Get admin from context
	data, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}
	admin, ok := data.(model.Admin)
	if !ok {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Invalid admin data"))
		return
	}

	// Get microsite ID from URL parameter
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	// Parse request body for comment (required for rejection)
	var reqBody struct {
		RejectionReason string `json:"rejection_reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Rejection comment is required"))
		return
	}

	// Fetch microsite
	var microsite model.MicroSite
	if err := database.DB.Where("uuid = ?", uuid).First(&microsite).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("Microsite not found"))
		return
	}

	// Update microsite status
	microsite.Status = "Rejected"
	microsite.AdminId = admin.Id
	microsite.RejectionReason = reqBody.RejectionReason

	if err := database.DB.Save(&microsite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to reject microsite"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Microsite rejected successfully",
		"data":    microsite,
	})
}

/*
* Delete any microsite (admin privilege)
* @param c *gin.Context
* @return gin.JSON
 */
func AdminDeleteMicrosite(c *gin.Context) {
	// Get admin from context
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Get microsite ID from URL parameter
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	// Check if microsite exists
	var microsite model.MicroSite
	if err := database.DB.Where("uuid = ?", id).First(&microsite).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("Microsite not found"))
		return
	}

	// Delete microsite (CASCADE will delete related Services and SocialLinks)
	if err := database.DB.Delete(&microsite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to delete microsite"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"message":    "Microsite deleted successfully",
		"deleted_id": id,
	})
}

// AdminGetMicrositeDetails - Get details of a specific microsite (admin privilege)
/*
* Get details of a specific microsite (admin privilege)
* @param c *gin.Context
* @return gin.JSON
 */
func AdminGetMicrositeDetails(c *gin.Context) {
	// Get admin from context
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Get microsite ID from URL parameter
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Microsite ID is required"))
		return
	}

	// Fetch microsite with related data
	var microsite model.MicroSite
	if err := database.DB.Where("uuid = ?", uuid).
		Preload("User").
		Preload("Services").
		Preload("SocialLinks").
		First(&microsite).Error; err != nil {
		c.JSON(http.StatusNotFound, service.ErrorResponse("Microsite not found"))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, service.SuccessResponse("Microsite details fetched successfully", microsite))
}
