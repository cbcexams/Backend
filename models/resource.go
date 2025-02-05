package models

import (
	"strings"
	"time"
)

// Resource represents an existing resource in the system
type Resource struct {
	Id                      string    `orm:"pk;column(id);type(uuid)" json:"id"`
	ParentUrl               string    `orm:"column(parent_url);type(text)" json:"parent_url"`
	GoogleDriveDownloadLink string    `orm:"column(google_drive_download_link);type(text)" json:"google_drive_download_link"`
	Name                    string    `orm:"column(name);type(text)" json:"name"`
	RelativePath            string    `orm:"column(relative_path);type(text)" json:"relative_path"`
	CreatedAt               time.Time `orm:"column(created_at);type(timestamp with time zone);auto_now_add" json:"created_at"`
	DjangoRelativePath      string    `orm:"column(django_relative_path);type(text)" json:"django_relative_path"`
	ParentDirectory         string    `orm:"column(parent_directory);type(text)" json:"parent_directory"`
	Categories              string    `orm:"column(categories);type(varchar)" json:"categories"`
}

// TableName specifies the database table name
func (r *Resource) TableName() string {
	return "web_crawler_resources"
}

// GetCategories returns the categories as a string slice
func (r *Resource) GetCategories() []string {
	// Remove the curly braces and split by comma
	categoriesStr := r.Categories[1 : len(r.Categories)-1] // Remove { and }
	if categoriesStr == "" {
		return []string{}
	}
	return strings.Split(categoriesStr, ",")
}

// SetCategories sets the categories from a string slice
func (r *Resource) SetCategories(categories []string) {
	r.Categories = "{" + strings.Join(categories, ",") + "}"
}
