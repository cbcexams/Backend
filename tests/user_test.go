package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"cbc-backend/context"
	"cbc-backend/utils"

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

func TestPasswordReset(t *testing.T) {
	// Create test user
	createTestUser(t)

	// Request password reset
	forgotBody := `{
		"email": "test@example.com"
	}`
	w := makeTestRequest(t, "POST", "/v1/user/forgot-password", forgotBody, "")
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	resetToken := data["reset_token"].(string)

	// Reset password
	resetBody := fmt.Sprintf(`{
		"reset_token": "%s",
		"new_password": "newpassword123"
	}`, resetToken)

	w = makeTestRequest(t, "POST", "/v1/user/reset-password", resetBody, "")
	assert.Equal(t, 200, w.Code)

	// Try login with new password
	loginBody := `{
		"username": "testuserN",
		"password": "newpassword123"
	}`
	w = makeTestRequest(t, "POST", "/v1/user/login", loginBody, "")
	assert.Equal(t, 200, w.Code)
}

func TestDeleteUser(t *testing.T) {
	// Create test user
	token := createTestUser(t)

	// Get user ID from token
	claims, err := utils.GetJWTClaims(&context.Context{Input: &context.BeegoInput{
		RequestBody: []byte{},
		Request:     &http.Request{Header: http.Header{"Authorization": []string{"Bearer " + token}}},
	}})
	assert.NoError(t, err)
	userID := fmt.Sprintf("%v", claims["user_id"])

	// Try to delete without token
	w := makeTestRequest(t, "DELETE", "/v1/user/"+userID, "", "")
	assert.Equal(t, 401, w.Code)

	// Try to delete with invalid token
	w = makeTestRequest(t, "DELETE", "/v1/user/"+userID, "", "invalid.token")
	assert.Equal(t, 401, w.Code)

	// Delete with valid token
	w = makeTestRequest(t, "DELETE", "/v1/user/"+userID, "", token)
	assert.Equal(t, 200, w.Code)

	// Verify user is deleted
	w = makeTestRequest(t, "POST", "/v1/user/login", `{
		"username": "testuserN",
		"password": "password12N3"
	}`, "")
	assert.Equal(t, 401, w.Code)
}
