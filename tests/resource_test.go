package test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func createTestFile(t *testing.T, filename string, content string) string {
	tmpDir := t.TempDir()
	filepath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filepath, []byte(content), 0666)
	if err != nil {
		t.Fatal(err)
	}
	return filepath
}

func createMultipartRequest(t *testing.T, filename string, fieldname string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("title", "Test Resource")
	_ = writer.WriteField("description", "Test Description")
	_ = writer.WriteField("level", "grade1-6")

	filepath := createTestFile(t, filename, "test content")
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/v1/resources", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestResourceUpload(t *testing.T) {
	Convey("Test Resource Upload", t, func() {
		Convey("Should accept PDF file", func() {
			req, err := createMultipartRequest(t, "test.pdf", "file")
			So(err, ShouldBeNil)

			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, req)

			logs.Info("testing", "TestResourceUpload PDF", "Code[%d]\n%s", w.Code, w.Body.String())
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should accept DOCX file", func() {
			req, err := createMultipartRequest(t, "test.docx", "file")
			So(err, ShouldBeNil)

			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, req)

			logs.Info("testing", "TestResourceUpload DOCX", "Code[%d]\n%s", w.Code, w.Body.String())
			So(w.Code, ShouldEqual, 200)
		})

		Convey("Should reject JPG file", func() {
			req, err := createMultipartRequest(t, "test.jpg", "file")
			So(err, ShouldBeNil)

			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, req)

			logs.Info("testing", "TestResourceUpload JPG", "Code[%d]\n%s", w.Code, w.Body.String())
			So(w.Body.String(), ShouldContainSubstring, "invalid file type")
		})

		Convey("Should reject MP4 file", func() {
			req, err := createMultipartRequest(t, "test.mp4", "file")
			So(err, ShouldBeNil)

			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, req)

			logs.Info("testing", "TestResourceUpload MP4", "Code[%d]\n%s", w.Code, w.Body.String())
			So(w.Body.String(), ShouldContainSubstring, "invalid file type")
		})
	})
}
