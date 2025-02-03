package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Resource represents a teaching resource in the system
// It maps to the web_crawler_resources table in the database
type Resource struct {
	Id                      string    `orm:"pk;column(id);type(uuid)" json:"id"`
	ParentUrl               string    `orm:"column(parent_url);type(text);null" json:"parent_url"`
	GoogleDriveDownloadLink string    `orm:"column(google_drive_download_link);type(text);unique;null" json:"google_drive_download_link"`
	Name                    string    `orm:"column(name);type(text)" json:"name"`
	RelativePath            string    `orm:"column(relative_path);type(text);unique" json:"relative_path"`
	CreatedAt               time.Time `orm:"column(created_at);type(timestamp with time zone);auto_now_add" json:"created_at"`
	DjangoRelativePath      string    `orm:"column(django_relative_path);type(text);unique;null" json:"django_relative_path"`
	ParentDirectory         string    `orm:"column(parent_directory);type(text);null" json:"parent_directory"`
	Categories              string    `orm:"column(categories);type(text);null" json:"categories"`
}

// TableName specifies the database table name for the Resource model
func (r *Resource) TableName() string {
	return "web_crawler_resources"
}

// GetCategories returns the categories as a slice
// If Categories is empty, returns nil
func (r *Resource) GetCategories() []string {
	if r.Categories == "" {
		return nil
	}
	return strings.Split(r.Categories, ",")
}

// SetCategories sets the categories from a slice
// If the slice is empty, sets Categories to empty string
func (r *Resource) SetCategories(categories []string) {
	if len(categories) == 0 {
		r.Categories = ""
		return
	}
	r.Categories = strings.Join(categories, ",")
}

// SearchResources searches resources with pagination
func SearchResources(params map[string]string, page int) (*Pagination, error) {
	var resources []*Resource
	o := orm.NewOrm()

	// Log search start
	fmt.Println("\n==================================================")
	fmt.Println("                Resource Search                    ")
	fmt.Println("==================================================")
	fmt.Printf("Search Parameters: %+v\n", params)
	fmt.Printf("Page Number: %d\n", page)

	// Build the base query with search conditions
	baseQuery := `
		SELECT id, parent_url, google_drive_download_link, name, 
			   relative_path, created_at, django_relative_path, 
			   parent_directory, categories
		FROM web_crawler_resources
		WHERE 1=1
	`
	var conditions []string
	var queryParams []interface{}

	// Add name search if provided
	if name, ok := params["name"]; ok && name != "" {
		// Make the search more precise by using word boundaries
		conditions = append(conditions, "(name ILIKE ? OR name ILIKE ? OR name ILIKE ?)")
		searchTerm := name
		queryParams = append(queryParams,
			searchTerm+"%",      // Starts with
			"% "+searchTerm+"%", // Contains as whole word
			"%"+searchTerm+"%",  // Contains anywhere
		)
	}

	// Add categories search if provided
	if categories, ok := params["categories"]; ok && categories != "" {
		conditions = append(conditions, "categories::text ILIKE ?")
		queryParams = append(queryParams, "%\""+categories+"\"%") // Look for exact category in JSON array
	}

	// Combine conditions
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	baseQuery += " ORDER BY created_at DESC"

	// Add pagination
	baseQuery += " LIMIT ? OFFSET ?"

	// Calculate pagination parameters
	const PageSize = 20
	offset := (page - 1) * PageSize
	queryParams = append(queryParams, PageSize, offset)

	// Execute count query
	fmt.Println("\n[1] Executing Count Query...")
	var total int64

	// Build count query
	countQuery := `
		SELECT COUNT(*)
		FROM web_crawler_resources
		WHERE 1=1
	`
	if len(conditions) > 0 {
		countQuery += " AND " + strings.Join(conditions, " AND ")
	}

	err := o.Raw(countQuery, queryParams[:len(queryParams)-2]...).QueryRow(&total)
	if err != nil {
		fmt.Printf("❌ Error counting resources: %v\n", err)
		return nil, fmt.Errorf("error counting resources: %v", err)
	}
	fmt.Printf("✅ Total records found: %d\n", total)

	// Execute main query
	fmt.Printf("\n[2] Executing Main Query...\n")
	fmt.Printf("Query: %s\n", baseQuery)
	fmt.Printf("Parameters: %+v\n", queryParams)

	num, err := o.Raw(baseQuery, queryParams...).QueryRows(&resources)
	if err != nil {
		fmt.Printf("❌ Error fetching resources: %v\n", err)
		return nil, fmt.Errorf("error fetching resources: %v", err)
	}
	fmt.Printf("✅ Retrieved %d records\n", num)

	// Calculate pagination information
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))

	// Return pagination result
	return &Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PageSize:    PageSize,
		Items:       resources,
	}, nil
}

// Pagination represents a paginated result set
type Pagination struct {
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	PageSize    int         `json:"page_size"`
	Items       []*Resource `json:"items"`
}
