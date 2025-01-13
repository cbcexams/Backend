package controllers

import (
	"cbc-backend/models"
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

type JobController struct {
	beego.Controller
}

// @Title Get
// @Description get jobs with pagination and filters
// @Param	title	query	string	false	"Job title search"
// @Param	type	query	string	false	"Job type (Full-time, Part-time, Contract)"
// @Param	location	query	string	false	"Job location"
// @Param	page	query	int	false	"Page number (default: 1)"
// @Success 200 {object} models.JobPagination
// @router / [get]
func (j *JobController) Get() {
	fmt.Println("Job Get method called")

	params := make(map[string]string)
	params["title"] = j.GetString("title")
	params["type"] = j.GetString("type")
	params["location"] = j.GetString("location")

	page, err := j.GetInt("page")
	if err != nil {
		page = 1
	}

	fmt.Printf("Search params: %v, Page: %d\n", params, page)

	result, err := models.SearchJobs(params, page)
	if err != nil {
		fmt.Printf("Error getting jobs: %v\n", err)
		j.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		fmt.Printf("Successfully retrieved %d jobs (Page %d of %d)\n",
			len(result.Items), result.CurrentPage, result.TotalPages)
		j.Data["json"] = result
	}
	j.ServeJSON()
}

// @Title Post
// @Description create new job posting
// @Param	body	body	models.Job	true	"Job content"
// @Success 200 {object} map[string]interface{}
// @Failure 403 body is empty
// @router / [post]
func (j *JobController) Post() {
	var job models.Job
	err := json.Unmarshal(j.Ctx.Input.RequestBody, &job)
	if err != nil {
		j.Data["json"] = map[string]string{"error": "Invalid request body"}
		j.ServeJSON()
		return
	}

	id, err := models.AddJob(job)
	if err != nil {
		j.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		j.Data["json"] = map[string]interface{}{
			"id":      id,
			"message": "Job created successfully",
		}
	}
	j.ServeJSON()
}
