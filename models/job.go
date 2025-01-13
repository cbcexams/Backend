package models

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Job struct {
	Id          int       `orm:"pk;auto" json:"id"`
	Title       string    `orm:"size(100)" json:"title"`
	Description string    `orm:"type(text)" json:"description"`
	Location    string    `orm:"size(100)" json:"location"`
	Type        string    `orm:"size(50)" json:"type"` // Full-time, Part-time, Contract
	Salary      string    `orm:"size(50)" json:"salary"`
	CreatedAt   time.Time `orm:"auto_now_add;type(timestamp)" json:"created_at"`
}

const JobPageSize = 20

type JobPagination struct {
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int64  `json:"total_items"`
	PageSize    int    `json:"page_size"`
	Items       []*Job `json:"items"`
}

func init() {
	orm.RegisterModel(new(Job))
}

func AddJob(job Job) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(&job)
	if err != nil {
		fmt.Printf("Error adding job: %v\n", err)
		return 0, err
	}
	return id, nil
}

func SearchJobs(params map[string]string, page int) (*JobPagination, error) {
	var jobs []*Job
	o := orm.NewOrm()

	query := o.QueryTable("job")

	// Apply filters
	for key, value := range params {
		if value != "" {
			switch key {
			case "title":
				query = query.Filter("title__icontains", value)
			case "type":
				query = query.Filter("type", value)
			case "location":
				query = query.Filter("location__icontains", value)
			}
		}
	}

	// Get total count
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	totalPages := int((total + int64(JobPageSize) - 1) / int64(JobPageSize))
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}
	offset := (page - 1) * JobPageSize

	// Get paginated results
	_, err = query.OrderBy("-created_at").Limit(JobPageSize, offset).All(&jobs)
	if err != nil {
		return nil, err
	}

	return &JobPagination{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PageSize:    JobPageSize,
		Items:       jobs,
	}, nil
}
