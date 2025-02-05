package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type JobController struct {
	beego.Controller
}

// Get retrieves a list of jobs
func (c *JobController) Get() {
	fmt.Println("\n==================================================")
	fmt.Println("              Jobs Search                          ")
	fmt.Println("==================================================")

	// Get all possible query parameters
	page, _ := c.GetInt("page", 1)
	params := map[string]string{
		"title":    c.GetString("title"),
		"type":     c.GetString("type"),
		"location": c.GetString("location"),
	}

	// Log search parameters
	fmt.Printf("Search Parameters:\n")
	for k, v := range params {
		if v != "" {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}
	fmt.Printf("- page: %d\n", page)

	// Ensure jobs table exists
	if err := models.EnsureJobsTable(); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to ensure jobs table exists", nil, err)
		return
	}

	pagination, err := models.SearchJobs(params, page)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to search jobs", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "Jobs retrieved successfully", pagination, nil)
}

// GetOne retrieves a single job by ID
func (c *JobController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	job := models.Job{Id: id}
	err = models.GetJob(&job)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "", job, nil)
}

// Post creates a new job
func (c *JobController) Post() {
	// Ensure jobs table exists
	if err := models.EnsureJobsTable(); err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to ensure jobs table exists", nil, err)
		return
	}

	var job models.Job
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &job); err != nil {
		utils.SendResponse(&c.Controller, false, "Invalid request body", nil, err)
		return
	}

	id, err := models.AddJob(job)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "Failed to create job", nil, err)
		return
	}

	job.Id = int(id)
	utils.SendResponse(&c.Controller, true, "Job created successfully", job, nil)
}

// Put updates an existing job
func (c *JobController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	var job models.Job
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &job); err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	job.Id = id
	err = models.UpdateJob(&job)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "Job updated successfully", job, nil)
}

// Delete removes a job
func (c *JobController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	err = models.DeleteJob(id)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "Job deleted successfully", nil, nil)
}

func init() {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		fmt.Printf("Failed to create uploads directory: %v\n", err)
	}
}
