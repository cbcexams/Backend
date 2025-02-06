package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompleteUserFlow(t *testing.T) {
	// 1. User signup
	token := createTestUser(t)

	// 2. Create job
	jobID := createTestJob(t, token)

	// 3. Upload resource
	resourceID := uploadTestResource(t, token)

	// 4. Verify everything
	w := makeTestRequest(t, "GET", "/v1/jobs/"+jobID, "", token)
	assert.Equal(t, 200, w.Code)

	w = makeTestRequest(t, "GET", "/v1/resources/"+resourceID, "", token)
	assert.Equal(t, 200, w.Code)
}
