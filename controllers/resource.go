package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

// ResourceController handles resource-related operations
type ResourceController struct {
	beego.Controller
}

// Helper methods for error responses
func (c *ResourceController) ServerError(msg string, err error) {
	utils.SendResponse(&c.Controller, false, msg, nil, err)
}

func (c *ResourceController) BadRequest(msg string, err error) {
	utils.SendResponse(&c.Controller, false, msg, nil, err)
}

// GetUserID gets the user ID from JWT token
func (c *ResourceController) GetUserID() int {
	authHeader := c.Ctx.Input.Header("Authorization")
	if authHeader == "" {
		return 0 // Or handle error as needed
	}

	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Get user ID from token
	userID, err := utils.GetUserIDFromToken(tokenString)
	if err != nil {
		return 0 // Or handle error as needed
	}

	return userID
}

// Get retrieves a list of resources with pagination
func (r *ResourceController) Get() {
	// Get all query parameters
	page, _ := r.GetInt("page", 1)
	pageSize := 20 // Fixed page size

	// Create params map with all search parameters
	params := map[string]string{
		"name": r.GetString("name"),
	}

	// Log the received parameters for debugging
	fmt.Printf("Received search parameters: %+v\n", params)

	// Fetch resources with pagination
	o := orm.NewOrm()
	var resources []*models.Resource

	// Calculate offset
	offset := (page - 1) * pageSize

	// Build query
	qs := o.QueryTable("web_crawler_resources")
	if name, ok := params["name"]; ok && name != "" {
		qs = qs.Filter("name__icontains", name)
	}

	// Get total count
	totalItems, err := qs.Count()
	if err != nil {
		utils.SendResponse(&r.Controller, false, "Failed to get total count", nil, err)
		return
	}

	// Calculate total pages
	totalPages := (totalItems + int64(pageSize) - 1) / int64(pageSize)

	// Get paginated results
	_, err = qs.OrderBy("-created_at").Limit(pageSize).Offset(offset).All(&resources)
	if err != nil {
		utils.SendResponse(&r.Controller, false, "Failed to fetch resources", nil, err)
		return
	}

	// Create pagination response with fields in the desired order
	pagination := map[string]interface{}{
		"current_page": page,
		"total_pages":  totalPages,
		"page_size":    pageSize,
		"total_items":  totalItems,
		"items":        resources, // items last
	}

	utils.SendResponse(&r.Controller, true, "", pagination, nil)
}

// Post handles file upload
func (c *ResourceController) Post() {
	// Ensure uploads table exists
	if err := models.EnsureUploadsTable(); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to ensure uploads table exists", nil, err)
		return
	}

	// Get the uploaded file
	file, header, err := c.GetFile("file")
	if err != nil {
		utils.SendResponse(&c.Controller, false, "No file uploaded", nil, err)
		return
	}
	defer file.Close()

	// Validate file type
	if err := validateFileType(header.Filename); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid file type", nil, err)
		return
	}

	// Generate unique ID for the upload
	uploadId := uuid.New().String()

	// Create upload directory if it doesn't exist
	uploadDir := "uploads/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to create upload directory", nil, err)
		return
	}

	// Generate unique filename
	ext := path.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uploadId, ext)
	filePath := path.Join(uploadDir, fileName)

	// Save the file
	if err := c.SaveToFile("file", filePath); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to save file", nil, err)
		return
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to get file info", nil, err)
		return
	}

	// Create upload record
	upload := &models.Upload{
		Id:          uploadId,
		FileName:    header.Filename,
		FilePath:    filePath,
		FileSize:    fileInfo.Size(),
		ContentType: header.Header.Get("Content-Type"),
		UserID:      c.GetUserID(),
	}

	// Save to database
	if err := models.CreateUpload(upload); err != nil {
		// Clean up the file if database insert fails
		os.Remove(filePath)
		utils.SendResponse(&c.Controller, false, "Failed to create upload record", nil, err)
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"upload_id": uploadId,
			"file_name": header.Filename,
			"file_path": filePath,
			"file_size": fileInfo.Size(),
		},
	}
	c.ServeJSON()
}

// validateFileType checks if the file extension is allowed
func validateFileType(filename string) error {
	ext := strings.ToLower(path.Ext(filename))
	validTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".txt":  true,
		".rtf":  true,
	}

	if !validTypes[ext] {
		return fmt.Errorf("invalid file type: %s", ext)
	}
	return nil
}

// Any handles unmatched GET requests for debugging
func (r *ResourceController) Any() {
	fmt.Printf("Path: %s\n", r.Ctx.Request.URL.Path)
	fmt.Printf("Method: %s\n", r.Ctx.Request.Method)
	utils.SendResponse(&r.Controller, true, "Debug route hit", nil, nil)
}
