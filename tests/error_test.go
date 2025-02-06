package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidRequests(t *testing.T) {
	token := createTestUser(t)

	// Test invalid JSON
	w := makeTestRequest(t, "POST", "/v1/jobs", "invalid json", token)
	assert.Equal(t, 400, w.Code)

	// Test missing required fields
	w = makeTestRequest(t, "POST", "/v1/jobs", "{}", token)
	assert.Equal(t, 400, w.Code)
}
