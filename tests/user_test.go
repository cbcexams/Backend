package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSignup(t *testing.T) {
	body := `{
		"username": "newuser",
		"password": "password123",
		"email": "new@example.com",
		"role": "teacher"
	}`

	w := makeTestRequest(t, "POST", "/v1/user/signup", body, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
}

func TestUserLogin(t *testing.T) {
	// First create a user
	createTestUser(t)

	// Try to login
	body := `{
		"username": "testuser",
		"password": "testpass"
	}`

	w := makeTestRequest(t, "POST", "/v1/user/login", body, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
}

func TestInvalidLogin(t *testing.T) {
	body := `{
		"username": "nonexistent",
		"password": "wrongpass"
	}`

	w := makeTestRequest(t, "POST", "/v1/user/login", body, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestUserFlow(t *testing.T) {
	// Test signup
	signupBody := `{
		"username": "testuser",
		"password": "password123",
		"email": "test@example.com",
		"role": "teacher"
	}`

	w := makeTestRequest(t, "POST", "/v1/user/signup", signupBody, "")
	assert.Equal(t, 200, w.Code)

	// Test login
	loginBody := `{
		"username": "testuser",
		"password": "password123"
	}`

	w = makeTestRequest(t, "POST", "/v1/user/login", loginBody, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	token := data["token"].(string)
	assert.NotEmpty(t, token)

	// Test logout
	w = makeTestRequest(t, "GET", "/v1/user/logout", "", token)
	assert.Equal(t, 200, w.Code)
}
