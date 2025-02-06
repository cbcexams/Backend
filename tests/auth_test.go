package tests

import (
	"cbc-backend/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTTokenValidation(t *testing.T) {
	token := createTestUser(t)

	// Test valid token
	claims, err := utils.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Test invalid token
	_, err = utils.ValidateJWT("invalid.token.here")
	assert.Error(t, err)
}

func TestProtectedRoutes(t *testing.T) {
	// Test without token
	w := makeTestRequest(t, "POST", "/v1/jobs", "{}", "")
	assert.Equal(t, 401, w.Code)

	// Test with invalid token
	w = makeTestRequest(t, "POST", "/v1/jobs", "{}", "invalid.token")
	assert.Equal(t, 401, w.Code)
}
