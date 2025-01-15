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
