package controllers

import (
	"cbc-backend/models"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var allowedFileTypes = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".txt":  true,
	".rtf":  true,
}

func validateFileType(filename string) error {
	ext := strings.ToLower(path.Ext(filename))
	if !allowedFileTypes[ext] {
		return errors.New("invalid file type: only PDF and document files are allowed")
	}
	return nil
}

// ResourceController operations about Resources
type ResourceController struct {
	beego.Controller
}

// @Title Get
// @Description get resources by level or search with pagination
// @Param	level	query	string	false	"Resource level"
// @Param	title	query	string	false	"Resource title"
// @Param	description	query	string	false	"Resource description"
// @Param	page	query	int	false	"Page number (default: 1)"
// @Success 200 {object} models.Pagination
// @Failure 403 :level is empty
// @router / [get]
func (r *ResourceController) Get() {
	fmt.Println("Resource Get method called")

	// Get all query parameters
	params := make(map[string]string)
	params["level"] = r.GetString("level")
	params["title"] = r.GetString("title")
	params["description"] = r.GetString("description")

	// Get page number, default to 1 if not provided
	page, err := r.GetInt("page")
	if err != nil {
		page = 1
	}

	fmt.Printf("Search params: %v, Page: %d\n", params, page)

	result, err := models.SearchResources(params, page)
	if err != nil {
		fmt.Printf("Error getting resources: %v\n", err)
		r.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		fmt.Printf("Successfully retrieved %d resources (Page %d of %d)\n",
			len(result.Items), result.CurrentPage, result.TotalPages)
		r.Data["json"] = result
	}
	r.ServeJSON()
}

// @Title Post
// @Description create resources
// @Param	title	formData	string	true	"Resource title"
// @Param	description	formData	string	true	"Resource description"
// @Param	level	formData	string	true	"Resource level"
// @Param	file	formData	file	true	"Resource file"
// @Success 200 {object} models.Resource
// @Failure 403 body is empty
// @router / [post]
func (r *ResourceController) Post() {
	title := r.GetString("title")
	description := r.GetString("description")
	level := r.GetString("level")

	f, h, err := r.GetFile("file")
	if err != nil {
		r.Data["json"] = map[string]string{"error": err.Error()}
		r.ServeJSON()
		return
	}
	defer f.Close()

	if err := validateFileType(h.Filename); err != nil {
		r.Data["json"] = map[string]string{"error": err.Error()}
		r.ServeJSON()
		return
	}

	ext := path.Ext(h.Filename)
	filename := time.Now().Format("20060102150405") + ext
	filepath := "uploads/" + filename

	err = r.SaveToFile("file", filepath)
	if err != nil {
		r.Data["json"] = map[string]string{"error": err.Error()}
		r.ServeJSON()
		return
	}

	resource := models.Resource{
		Title:       title,
		Description: description,
		FilePath:    filepath,
		Level:       level,
	}

	id, err := models.AddResource(resource)
	if err != nil {
		r.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		r.Data["json"] = map[string]interface{}{
			"id":      id,
			"message": "Resource created successfully",
		}
	}
	r.ServeJSON()
}
