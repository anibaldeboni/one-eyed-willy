package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePdfCaseSuccess(t *testing.T) {
	setup()
	var (
		reqJSON = `{"html":"PGh0bWw+CjxoZWFkPgoJPHRpdGxlPk15IFBERiBGaWxlPC90aXRsZT4KPC9oZWFkPgo8Ym9keT4KCTxwPkhlbGxvIHRoZXJlISBJJ20gYSBwZGYgZmlsZSBnZW5lcmF0ZSBmcm9tIGEgaHRtbCB1c2luZyBnbyBhbmQgZ29wZGYgcGFja2FnZTwvcD4KPC9ib2R5Pgo8L2h0bWw+"}`
	)
	req := httptest.NewRequest(echo.POST, "/pdf", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GeneratePdfFromHTML(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "application/pdf", rec.Header().Clone().Get("Content-Type"))
		assert.Greater(t, len(rec.Body.Bytes()), 0)
		assert.NotEmpty(t, rec.Body)
	}
}

func TestMergePdfsCaseSuccess(t *testing.T) {
	setup()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file1, file2 := loadFiles(t)
	part, _ := writer.CreateFormFile("files", "file1.pdf")
	part.Write(file1)

	part2, _ := writer.CreateFormFile("files", "file2.pdf")
	part2.Write(file2)
	writer.Close()

	req := httptest.NewRequest(echo.POST, "/pdf/merge", body)
	contentType := fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, writer.Boundary())
	req.Header.Add(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NotPanics(t, func() { _ = h.MergePdfs(c) }) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NoError(t, h.MergePdfs(c))
	}
}

func loadFiles(t *testing.T) (file1 []byte, file2 []byte) {
	file1, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	file2, err = os.ReadFile("../../testdata/file2.pdf")
	if err != nil {
		t.Fatal("testdata/file2.pdf not found")
	}

	return file1, file2
}
