package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
