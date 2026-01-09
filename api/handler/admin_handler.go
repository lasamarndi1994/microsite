package handler

import (
	"encoding/csv"
	"fmt"
	"math"
	"micro-site/api/model"
	"micro-site/database"
	"micro-site/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

	// Filter: Only users with at least one microsite
	query = query.Where("EXISTS (SELECT 1 FROM micro_sites WHERE micro_sites.user_id = users.id)")

	// Search filter
	search := c.Query("search")
	if search != "" {
		searchLike := "%" + search + "%"
		query = query.Where("mobile_number LIKE ? OR email LIKE ? OR CAST(mobile_number AS CHAR) LIKE ?", searchLike, searchLike, searchLike)
	}

	// Apply status filter if provided

	// Date filter
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	if fromDate != "" && toDate != "" {
		query = query.Where("DATE(updated_at) BETWEEN ? AND ?", fromDate, toDate)
	} else if fromDate != "" {
		query = query.Where("DATE(updated_at) >= ?", fromDate)
	} else if toDate != "" {
		query = query.Where("DATE(updated_at) <= ?", toDate)
	}

	// Count total records
	var totalRecords int64
	query.Count(&totalRecords)

	// Execute query with pagination
	var users []model.User
	// Sort by latest microsite creation time, then by user updated_at
	if err := query.Order("(SELECT MAX(id) FROM micro_sites WHERE micro_sites.user_id = users.id) DESC").Order("updated_at desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch users"))
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	// Prepare response with masked mobile numbers
	type UserResponse struct {
		model.User
		MobileNumber string `json:"mobile_number"`
	}
	var responseData []UserResponse

	// Get microsite counts for fetched users
	if len(users) > 0 {
		var userIDs []uint64
		for _, u := range users {
			userIDs = append(userIDs, u.Id)
		}

		type UserCount struct {
			UserId        int64
			TotalCount    int64
			RejectedCount int64
			ApprovedCount int64
			PendingCount  int64
		}
		var counts []UserCount

		// Aggregate counts
		database.DB.Model(&model.MicroSite{}).
			Select(`
				user_id, 
				count(*) as total_count,
				sum(case when status = 'Rejected' then 1 else 0 end) as rejected_count,
				sum(case when status = 'Approved' then 1 else 0 end) as approved_count,
				sum(case when status = 'Pending' then 1 else 0 end) as pending_count
			`).
			Where("user_id IN ?", userIDs).
			Group("user_id").
			Scan(&counts)

		// Map counts to users
		countMap := make(map[uint64]UserCount)
		for _, c := range counts {
			countMap[uint64(c.UserId)] = c
		}

		for _, u := range users {
			c := countMap[u.Id]
			u.MicrositeCount = c.TotalCount
			u.MicrositeRejectedCount = c.RejectedCount
			u.MicrositeApprovedCount = c.ApprovedCount
			u.MicrositePendingCount = c.PendingCount

			// Mask mobile number
			mobileStr := strconv.Itoa(u.MobileNumber)
			maskedMobile := mobileStr
			if len(mobileStr) > 5 {
				maskedMobile = "*****" + mobileStr[5:]
			}

			responseData = append(responseData, UserResponse{
				User:         u,
				MobileNumber: maskedMobile,
			})
		}
	} else {
		responseData = []UserResponse{}
	}

	// Return success response with pagination
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Users fetched successfully",
		"data":    responseData,
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

	// Count approved microsites
	var pendingCount int64
	if err := database.DB.Model(&model.MicroSite{}).Where("user_id = ? AND status = ?", user.Id, "Pending").Count(&pendingCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to count approved microsites"))
		return
	}

	// Mask mobile number
	mobileStr := strconv.Itoa(user.MobileNumber)
	maskedMobile := mobileStr
	if len(mobileStr) > 5 {
		maskedMobile = "*****" + mobileStr[5:]
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
			"mobile_number": maskedMobile,
			"status":        user.Status,
			"user_avatar":   user.UserAvatar,
			"user_slug":     user.Slug,
		},
		"data":          microsites,
		"count":         len(microsites),
		"pending_count": pendingCount,
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

/*
* Upload user CSV
* @param c *gin.Context
* @return gin.JSON
 */
func UploadUserCSV(c *gin.Context) {
	fmt.Println("UploadUserCSV")
	// Get admin from context
	_, exists := c.Get("admin")
	if !exists {
		c.JSON(http.StatusUnauthorized, service.ErrorResponse("Admin authentication required"))
		return
	}

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("File is required"))
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to open file"))
		return
	}
	defer src.Close()

	// Parse CSV
	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("Failed to parse CSV file"))
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, service.ErrorResponse("CSV file is empty or missing header"))
		return
	}

	// Map headers
	header := records[0]
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.TrimSpace(h)] = i
	}

	// Validate required columns
	requiredColumns := []string{"RemeshireCode", "RemeshireName", "EmailId", "MobileNo"}
	for _, col := range requiredColumns {
		if _, ok := headerMap[col]; !ok {
			c.JSON(http.StatusBadRequest, service.ErrorResponse("Missing required column: "+col))
			return
		}
	}
	fmt.Println("UploadUserCSV1")

	// Fetch existing mobile numbers
	var existingMobileNumbers []int
	if err := database.DB.Model(&model.User{}).Pluck("mobile_number", &existingMobileNumbers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch existing mobile numbers"))
		return
	}

	// Create a map for faster lookup
	existingMobileMap := make(map[int]bool)
	for _, num := range existingMobileNumbers {
		existingMobileMap[num] = true
	}

	// Fetch existing emails
	var existingEmails []string
	if err := database.DB.Model(&model.User{}).Pluck("email", &existingEmails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch existing emails"))
		return
	}

	// Create a map for faster lookup
	existingEmailMap := make(map[string]bool)
	for _, email := range existingEmails {
		existingEmailMap[email] = true
	}

	// Fetch existing slugs
	var existingSlugs []string
	if err := database.DB.Model(&model.User{}).Pluck("slug", &existingSlugs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to fetch existing slugs"))
		return
	}

	// Create a map for faster lookup
	existingSlugMap := make(map[string]bool)
	for _, s := range existingSlugs {
		existingSlugMap[s] = true
	}

	var users []model.User
	var errors []string
	var skippedCount int

	// Process records
	for i, record := range records[1:] {
		rowNum := i + 2 // 1-based index, +1 for header

		remeshireCode := record[headerMap["RemeshireCode"]]
		remeshireName := record[headerMap["RemeshireName"]]
		emailId := strings.ToLower(strings.TrimSpace(record[headerMap["EmailId"]]))
		mobileNoStr := record[headerMap["MobileNo"]]

		if remeshireCode == "" || remeshireName == "" || emailId == "" || mobileNoStr == "" {
			errors = append(errors, "Row "+strconv.Itoa(rowNum)+": Missing required fields")
			continue
		}

		mobileNo, err := strconv.Atoi(mobileNoStr)
		if err != nil {
			errors = append(errors, "Row "+strconv.Itoa(rowNum)+": Invalid MobileNo")
			continue
		}

		// Check if mobile number already exists
		if existingMobileMap[mobileNo] {
			skippedCount++
			continue
		}

		// Check if email already exists
		if existingEmailMap[emailId] {
			skippedCount++
			continue
		}

		// Check if slug already exists
		userSlug := slug.Make(remeshireName)
		if existingSlugMap[userSlug] {
			skippedCount++
			continue
		}

		fmt.Println("UploadUserCSV3")
		user := model.User{
			UserCode:     remeshireCode,
			UserName:     remeshireName,
			Email:        emailId,
			MobileNumber: mobileNo,
			Slug:         userSlug,
			Status:       "Active", // Default status
		}
		users = append(users, user)

		// Add to map to handle duplicates within the CSV itself
		existingMobileMap[mobileNo] = true
		existingEmailMap[emailId] = true
		existingSlugMap[userSlug] = true
	}

	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "No valid records found to insert",
			"errors":  errors,
			"skipped": skippedCount,
		})
		return
	}

	fmt.Println("UploadUserCSV4")
	// Batch insert
	batchSize := 100
	if err := database.DB.CreateInBatches(users, batchSize).Error; err != nil {
		// If batch insert fails, we might want to return more specific errors, but for now generic
		c.JSON(http.StatusInternalServerError, service.ErrorResponse("Failed to insert users: "+err.Error()))
		return
	}

	response := gin.H{
		"status":        true,
		"message":       "Users uploaded successfully",
		"total_records": len(records) - 1,
		"inserted":      len(users),
		"skipped":       skippedCount,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}
