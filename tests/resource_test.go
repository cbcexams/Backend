package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

func TestResourceUpload(t *testing.T) {
	token := createTestUser(t)

	// Create test file
	content := []byte("test content")
	tmpfile, err := os.CreateTemp("", "test*.pdf")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.Write(content)
	assert.NoError(t, err)

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(tmpfile.Name()))
	assert.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader(content))
	assert.NoError(t, err)

	// Add other fields
	writer.WriteField("name", "Test Resource")
	writer.WriteField("categories", "math,test")
	writer.Close()

	// Make request
	r, _ := http.NewRequest("POST", "/v1/resources", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestResourceList(t *testing.T) {
	token := createTestUser(t)

	w := makeTestRequest(t, "GET", "/v1/resources", "", token)
	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "items")
}
