package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"fmt"
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
	// Get form parameters
	name := r.GetString("name")
	parentDir := r.GetString("parent_directory")
	categories := strings.Split(r.GetString("categories"), ",")
	if len(categories) == 1 && categories[0] == "" {
		categories = nil
	}

	fmt.Printf("Creating new resource with Name: %s, Parent Directory: %s, Categories: %v\n", name, parentDir, categories)

	// Create new resource instance
	resource := models.Resource{
		Id:              uuid.New().String(),
		Name:            name,
		ParentDirectory: parentDir,
		CreatedAt:       time.Now(),
	}
	resource.SetCategories(categories)

	// Handle file upload
	f, h, err := r.GetFile("file")
	if err != nil {
		fmt.Printf("Error getting file: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}
	defer f.Close()

	// Validate file type
	if err := validateFileType(h.Filename); err != nil {
		fmt.Printf("Invalid file type: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}

	// Generate file paths
	ext := path.Ext(h.Filename)
	timestamp := time.Now().Format("20060102150405")
	relativePath := fmt.Sprintf("uploads/%s%s", timestamp, ext)
	djangoPath := fmt.Sprintf("/media/resources/%s%s", timestamp, ext)

	resource.RelativePath = relativePath
	resource.DjangoRelativePath = djangoPath

	// Save uploaded file
	fmt.Printf("Saving file to path: %s\n", relativePath)
	err = r.SaveToFile("file", relativePath)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}

	// Save to database
	o := orm.NewOrm()
	fmt.Printf("Inserting resource into database...\n")
	_, err = o.Insert(&resource)
	if err != nil {
		fmt.Printf("Error inserting resource: %v\n", err)
		utils.SendResponse(&r.Controller, false, "", nil, err)
		return
	}

	// Send successful response
	fmt.Printf("Resource created successfully: %+v\n", resource)
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
