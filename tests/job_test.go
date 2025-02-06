package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobCRUD(t *testing.T) {
	token := createTestUser(t)

	// Test create job
	jobBody := `{
		"title": "Math Teacher",
		"description": "Looking for experienced math teacher",
		"location": "Nairobi",
		"type": "Full-time",
		"salary": "50000-70000"
	}`

	w := makeTestRequest(t, "POST", "/v1/jobs", jobBody, token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	jobID := int(data["id"].(float64))

	// Test get job
	w = makeTestRequest(t, "GET", fmt.Sprintf("/v1/jobs/%d", jobID), "", token)
	assert.Equal(t, 200, w.Code)

	// Test update job
	updateBody := `{
		"title": "Senior Math Teacher",
		"salary": "60000-80000"
	}`
	w = makeTestRequest(t, "PUT", fmt.Sprintf("/v1/jobs/%d", jobID), updateBody, token)
	assert.Equal(t, 200, w.Code)

	// Test search jobs
	w = makeTestRequest(t, "GET", "/v1/jobs?title=math&type=Full-time", "", token)
	assert.Equal(t, 200, w.Code)

	// Test delete job
	w = makeTestRequest(t, "DELETE", fmt.Sprintf("/v1/jobs/%d", jobID), "", token)
	assert.Equal(t, 200, w.Code)
}

func TestCreateJob(t *testing.T) {
	token := createTestUser(t)

	body := `{
		"title": "Math Teacher",
		"description": "Looking for experienced math teacher",
		"location": "Nairobi",
		"type": "Full-time",
		"salary": "50000-70000"
	}`

	w := makeTestRequest(t, "POST", "/v1/jobs", body, token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "id")
}

func TestListJobs(t *testing.T) {
	token := createTestUser(t)

	w := makeTestRequest(t, "GET", "/v1/jobs", "", token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "items")
}

func TestSearchJobs(t *testing.T) {
	token := createTestUser(t)

	w := makeTestRequest(t, "GET", "/v1/jobs?title=Math&type=Full-time", "", token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "items")
}
