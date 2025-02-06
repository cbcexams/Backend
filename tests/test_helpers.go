package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

const testDBConn = "user=postgres password=0000 dbname=cbcexams_test sslmode=disable"

// makeTestRequest is a helper function to make HTTP requests in tests
func makeTestRequest(t *testing.T, method, path, body, token string) *httptest.ResponseRecorder {
	r, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	assert.NoError(t, err)

	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}

// createTestUser creates a test user and returns the JWT token
func createTestUser(t *testing.T) string {
	body := `{
		"username": "testuserN",
		"password": "password12N3",
		"email": "test@example.com",
		"role": "teacher"
	}`

	w := makeTestRequest(t, "POST", "/v1/user/signup", body, "")
	assert.Equal(t, 200, w.Code)

	loginBody := `{
		"username": "testuserN",
		"password": "password12N3"
	}`

	w = makeTestRequest(t, "POST", "/v1/user/login", loginBody, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	return data["token"].(string)
}

//nolint:unused
func createTestJob(t *testing.T, token string) string {
	jobBody := `{
		"title": "Test Job",
		"description": "Test Description",
		"location": "Test Location",
		"type": "Full-time",
		"salary": "50000-70000"
	}`

	w := makeTestRequest(t, "POST", "/v1/jobs", jobBody, token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data := response["data"].(map[string]interface{})
	return fmt.Sprintf("%d", int(data["id"].(float64)))
}

// uploadTestResource creates a test resource and returns its ID
func uploadTestResource(t *testing.T, token string) string {
	// Implementation similar to TestResourceUpload
	// Returns the resource ID
	return "test-resource-id"
}
