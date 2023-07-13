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
	var (
		reqJSON = `{"html":"PGh0bWw+CjxoZWFkPgoJPHRpdGxlPk15IFBERiBGaWxlPC90aXRsZT4KPC9oZWFkPgo8Ym9keT4KCTxwPkhlbGxvIHRoZXJlISBJJ20gYSBwZGYgZmlsZSBnZW5lcmF0ZSBmcm9tIGEgaHRtbCB1c2luZyBnbyBhbmQgZ29wZGYgcGFja2FnZTwvcD4KPC9ib2R5Pgo8L2h0bWw+"}`
	)
	ctx, rec := setupServer(
		httptest.NewRequest(echo.POST, "/pdf", strings.NewReader(reqJSON)),
		echo.MIMEApplicationJSON,
	)

	assert.NoError(t, h.GeneratePdfFromHTML(ctx))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, MIMEApplicationPdf, rec.Header().Clone().Get("Content-Type"))
		assert.Greater(t, len(rec.Body.Bytes()), 0)
		assert.NotEmpty(t, rec.Body)
	}
}

func TestMergePdfs(t *testing.T) {
	tests := []struct {
		name          string
		numberOfFiles int
		httpStatus    int
	}{
		{
			name:          "when two files are provided",
			numberOfFiles: 2,
			httpStatus:    http.StatusOK,
		},
		{
			name:          "when just one file is provided",
			numberOfFiles: 1,
			httpStatus:    http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, boundary := createBody(t, tt.numberOfFiles)

			ctx, rec := setupServer(
				httptest.NewRequest(echo.POST, "/pdf/merge", body),
				fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary),
			)

			if assert.NotPanics(t, func() { _ = h.MergePdfs(ctx) }) {
				assert.Equal(t, tt.httpStatus, rec.Code)
				assert.NoError(t, h.MergePdfs(ctx))
			}
		})
	}
}

func createBody(t *testing.T, numberOfFiles int) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := os.ReadFile("../../testdata/file1.pdf")
	if err != nil {
		t.Fatal("testdata/file1.pdf not found")
	}

	for i := 0; i < numberOfFiles; i++ {
		part, _ := writer.CreateFormFile("files", fmt.Sprintf("file%d.pdf", i))
		if _, err := part.Write(file); err != nil {
			t.Fatal("could not write form-data")
		}
	}
	writer.Close()

	return body, writer.Boundary()
}

func setupServer(
	req *http.Request,
	contentType string,
) (echo.Context, *httptest.ResponseRecorder) {
	setup()
	req.Header.Set(echo.HeaderContentType, contentType)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
