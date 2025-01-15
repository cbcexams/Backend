package tests

import (
	"cbc-backend/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set test mode
	beego.BConfig.RunMode = "test"

	// Use test database
	testDBConn := "user=postgres password=0000 dbname=cbcexams_test sslmode=disable"
	models.InitDB(testDBConn)
}

// createTestUser creates a test user and returns the JWT token
func createTestUser(t *testing.T) string {
	user := models.User{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
		Role:     "teacher",
	}

	err := models.AddUser(&user)
	assert.NoError(t, err)

	token, err := user.GenerateToken()
	assert.NoError(t, err)

	return token
}

// makeTestRequest is a helper function to make HTTP requests in tests
func makeTestRequest(t *testing.T, method, path, body, token string) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	return w
}
