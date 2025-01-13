package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Resource struct {
	Id          int       `orm:"pk;auto"`
	Title       string    `orm:"size(100)"`
	Description string    `orm:"type(text)"`
	FilePath    string    `orm:"size(255)"`
	Level       string    `orm:"size(20)"`
	CreatedAt   time.Time `orm:"auto_now_add;type(timestamp)"`
}

func AddResource(r Resource) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(&r)
	return id, err
}

func GetResourcesByLevel(level string) ([]*Resource, error) {
	var resources []*Resource
	o := orm.NewOrm()

	query := o.QueryTable("resource")

	// If level is provided, filter by it
	if level != "" {
		query = query.Filter("level", level)
	}

	// Order by created_at descending (newest first)
	_, err := query.OrderBy("-created_at").All(&resources)

	if err != nil {
		fmt.Printf("Database query error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Found %d resources\n", len(resources))
	return resources, nil
}

const PageSize = 20

type Pagination struct {
	CurrentPage int         `json:"current_page"`
	TotalPages  int         `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	PageSize    int         `json:"page_size"`
	Items       []*Resource `json:"items"`
}

func SearchResources(params map[string]string, page int) (*Pagination, error) {
	var resources []*Resource
	o := orm.NewOrm()

	query := o.QueryTable("resource")

	// Apply filters based on provided parameters
	for key, value := range params {
		if value != "" {
			switch key {
			case "title":
				query = query.Filter("title__icontains", value)
			case "level":
				query = query.Filter("level", value)
			case "description":
				query = query.Filter("description__icontains", value)
			}
		}
	}

	// Get total count for pagination
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// Calculate pagination values
	totalPages := int((total + int64(PageSize) - 1) / int64(PageSize))
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	offset := (page - 1) * PageSize

	// Get paginated results
	_, err = query.OrderBy("-created_at").Limit(PageSize, offset).All(&resources)
	if err != nil {
		return nil, err
	}

	return &Pagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PageSize:    PageSize,
		Items:       resources,
	}, nil
}
