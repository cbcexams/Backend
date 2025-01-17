package controllers

import (
	"cbc-backend/models"
	"cbc-backend/utils"
	"encoding/json"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type JobController struct {
	beego.Controller
}

// Get retrieves a list of jobs
func (c *JobController) Get() {
	page, _ := c.GetInt("page", 1)
	title := c.GetString("title")
	jobType := c.GetString("type")
	location := c.GetString("location")

	params := map[string]string{
		"title":    title,
		"type":     jobType,
		"location": location,
	}

	pagination, err := models.SearchJobs(params, page)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	utils.SendResponse(&c.Controller, true, "", pagination, nil)
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
	var job models.Job
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &job); err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

	_, err := models.AddJob(job)
	if err != nil {
		utils.SendResponse(&c.Controller, false, "", nil, err)
		return
	}

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
