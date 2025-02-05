package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"fmt"
	"net/http"
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

// Get retrieves a list of resources with pagination
func (r *ResourceController) Get() {
	// Get all query parameters
	page, _ := r.GetInt("page", 1)

	// Create params map with all search parameters
	params := map[string]string{
		"name":       r.GetString("name"),
		"categories": r.GetString("categories"),
	}

	// Log the received parameters for debugging
	fmt.Printf("Received search parameters: %+v\n", params)

	// Fetch resources with pagination
	pagination, err := models.SearchResources(params, page)
	if err != nil {
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&r.Controller, true, "", pagination, nil)
}

// Post handles resource creation/upload
func (r *ResourceController) Post() {
	fmt.Println("\n==================================================")
	fmt.Println("              Resource Upload Started               ")
	fmt.Println("==================================================")

	// Log all form data received
	fmt.Println("\n[1] Received Form Data:")
	fmt.Printf("Name: %s\n", r.GetString("name"))
	fmt.Printf("Parent Directory: %s\n", r.GetString("parent_directory"))
	fmt.Printf("Categories: %s\n", r.GetString("categories"))

	// Get form parameters with validation
	name := r.GetString("name")
	if name == "" {
		fmt.Println("❌ Error: name is required")
		utils.SendResponse(&r.Controller, false, "", nil, fmt.Errorf("name is required"))
		return
	}

	parentDir := r.GetString("parent_directory")
	categoriesStr := r.GetString("categories")
	fmt.Printf("Raw categories string: %s\n", categoriesStr)

	// Clean up the categories string
	categoriesStr = strings.Trim(categoriesStr, "{}")
	var categories []string
	if categoriesStr != "" {
		categories = strings.Split(categoriesStr, ",")
		// Trim spaces from each category
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}
	}

	fmt.Println("\n[2] Parsed Parameters:")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Parent Directory: %s\n", parentDir)
	fmt.Printf("Categories: %v\n", categories)

	// Check if file was sent
	fmt.Println("\n[3] Processing File Upload:")
	f, h, err := r.GetFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			fmt.Println("❌ Error: No file was uploaded")
			utils.SendResponse(&r.Controller, false, "", nil, fmt.Errorf("file is required"))
			return
		}
		fmt.Printf("❌ Error getting file: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}
	defer f.Close()

	fmt.Printf("File received: %s\n", h.Filename)
	fmt.Printf("File size: %d bytes\n", h.Size)
	fmt.Printf("File header: %+v\n", h.Header)

	// Validate file type
	fmt.Println("\n[4] Validating File Type:")
	if err := validateFileType(h.Filename); err != nil {
		fmt.Printf("❌ Invalid file type: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}
	fmt.Println("✅ File type validation passed")

	// Generate file paths
	fmt.Println("\n[5] Generating File Paths:")
	ext := path.Ext(h.Filename)
	timestamp := time.Now().Format("20060102150405")
	relativePath := fmt.Sprintf("uploads/%s%s", timestamp, ext)
	djangoPath := fmt.Sprintf("/media/resources/%s%s", timestamp, ext)

	fmt.Printf("File extension: %s\n", ext)
	fmt.Printf("Timestamp: %s\n", timestamp)
	fmt.Printf("Relative path: %s\n", relativePath)
	fmt.Printf("Django path: %s\n", djangoPath)

	// Get user ID from JWT context
	fmt.Println("\n[6] Getting User ID from JWT:")
	userIDRaw := r.Ctx.Input.GetData("user_id")
	fmt.Printf("Raw user ID from JWT: %v (type: %T)\n", userIDRaw, userIDRaw)
	userID := userIDRaw.(float64)
	fmt.Printf("Converted user ID: %d\n", int(userID))

	// Create resource
	fmt.Println("\n[7] Creating Resource Object:")
	resource := models.Resource{
		Id:                 uuid.New().String(),
		UserID:             int(userID),
		Name:               name,
		ParentDirectory:    parentDir,
		CreatedAt:          time.Now(),
		RelativePath:       relativePath,
		DjangoRelativePath: djangoPath,
	}

	// Parse categories
	categoriesStr = r.GetString("categories")
	fmt.Printf("Raw categories string: %s\n", categoriesStr)

	// Clean up the categories string and set it
	categoriesStr = strings.Trim(categoriesStr, "{}")
	if categoriesStr != "" {
		categories := strings.Split(categoriesStr, ",")
		// Trim spaces from each category
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}
		resource.SetCategories(categories)
	}

	fmt.Printf("Resource object created: %+v\n", resource)

	// Save uploaded file
	fmt.Println("\n[8] Saving File to Disk:")
	fmt.Printf("Saving to path: %s\n", relativePath)
	if err := r.SaveToFile("file", relativePath); err != nil {
		fmt.Printf("❌ Error saving file: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}
	fmt.Println("✅ File saved successfully")

	// Save to database
	fmt.Println("\n[9] Saving to Database:")
	o := orm.NewOrm()
	fmt.Println("Creating database transaction...")

	// Debug: Print the SQL that will be executed
	sql := `
		INSERT INTO web_crawler_resources 
		(id, user_id, name, parent_directory, relative_path, django_relative_path, categories, google_drive_download_link) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL)
	`
	fmt.Printf("SQL to be executed: %s\n", sql)
	fmt.Printf("Values: id=%s, user_id=%d, name=%s, parent_dir=%s, rel_path=%s, django_path=%s, categories=%v\n",
		resource.Id, resource.UserID, resource.Name, resource.ParentDirectory,
		resource.RelativePath, resource.DjangoRelativePath, resource.Categories)

	_, err = o.Raw(sql, resource.Id, resource.UserID, resource.Name, resource.ParentDirectory,
		resource.RelativePath, resource.DjangoRelativePath, resource.Categories).Exec()
	if err != nil {
		fmt.Printf("❌ Database error: %v\n", err)
		// If file was saved but database insert failed, try to clean up the file
		fmt.Printf("Attempting to clean up saved file: %s\n", relativePath)
		if err := os.Remove(relativePath); err != nil {
			fmt.Printf("⚠️ Warning: Could not clean up file: %v\n", err)
		} else {
			fmt.Println("✅ File cleanup successful")
		}
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}
	fmt.Println("✅ Resource saved to database successfully")

	fmt.Println("\n[10] Upload Complete!")
	fmt.Printf("Resource ID: %s\n", resource.Id)
	fmt.Println("==================================================")

	utils.SendResponse(&r.Controller, true, "Resource created successfully", resource, nil)
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
