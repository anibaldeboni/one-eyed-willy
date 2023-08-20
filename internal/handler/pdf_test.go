package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/pkg/pdf"
	"github.com/one-eyed-willy/testdata"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePdfFromHTML(t *testing.T) {
	encodedHTML := base64.StdEncoding.EncodeToString([]byte("<h1>Hello World</h1>"))
	tests := []struct {
		name               string
		html               string
		pdfRender          *pdf.PdfRender
		wantHttpStatusCode int
		wantContentType    string
		wantErr            bool
	}{
		{
			name:               "when html is valid",
			html:               encodedHTML,
			pdfRender:          nil,
			wantHttpStatusCode: http.StatusOK,
			wantContentType:    MIMEApplicationPdf,
			wantErr:            false,
		},
		{
			name:               "when html is empty",
			html:               "",
			pdfRender:          nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
		{
			name:               "when html is not base64 encoded",
			html:               "this-is-not-a-base64-encoded-string",
			pdfRender:          nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
		{
			name:               "when PdfRender returns an error",
			html:               encodedHTML,
			pdfRender:          pdf.NewMockPdfRender(),
			wantHttpStatusCode: http.StatusInternalServerError,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		h, e := setupTestHandler(nil, tt.pdfRender)
		req := httptest.NewRequest(echo.POST, "/pdf/generate", strings.NewReader(fmt.Sprintf(`{"html":"%s"}`, tt.html)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(tt.html)))
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		t.Run(tt.name, func(t *testing.T) {
			if assert.NotPanics(t, func() { _ = h.GeneratePdfFromHTML(ctx) }) {
				assert.Equal(t, tt.wantHttpStatusCode, rec.Code)
				assert.Equal(t, tt.wantContentType, rec.Header().Get(echo.HeaderContentType))
				assert.NotEmpty(t, rec.Body)
			}
		})
	}
}

func BenchmarkGeneratePdf(b *testing.B) {
	h, e := setupTestHandler(nil, nil)
	req := httptest.NewRequest(echo.POST, "/pdf/generate", strings.NewReader(fmt.Sprintf(`{"html":"%s"}`, base64.StdEncoding.EncodeToString([]byte("<h1>Hello World</h1>")))))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len([]byte("<h1>Hello World</h1>"))))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	for i := 0; i < b.N; i++ {
		_ = h.GeneratePdfFromHTML(ctx)
	}
}

func TestMergePdfFiles(t *testing.T) {
	tests := []struct {
		name               string
		files              [][]byte
		pdfTool            *pdf.PdfTool
		wantHttpStatusCode int
		wantContentType    string
		wantErr            bool
	}{
		{
			name:               "when two files are provided",
			files:              [][]byte{testdata.LoadFile(t), testdata.LoadFile(t)},
			pdfTool:            nil,
			wantHttpStatusCode: http.StatusOK,
			wantContentType:    MIMEApplicationPdf,
			wantErr:            false,
		},
		{
			name:               "when just one file is provided",
			files:              [][]byte{testdata.LoadFile(t)},
			pdfTool:            nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
		{
			name:               "when the files can't be merged",
			files:              [][]byte{testdata.UnreadableFile(), testdata.UnreadableFile()},
			pdfTool:            nil,
			wantHttpStatusCode: http.StatusInternalServerError,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		h, e := setupTestHandler(tt.pdfTool, nil)
		body, boundary := createForm(t, tt.files)
		req := httptest.NewRequest(echo.POST, "/pdf/merge", body)
		req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
		req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		t.Run(tt.name, func(t *testing.T) {
			if assert.NotPanics(t, func() { _ = h.MergePdfs(ctx) }) {
				assert.Equal(t, tt.wantHttpStatusCode, rec.Code)
				assert.Equal(t, tt.wantContentType, rec.Header().Get(echo.HeaderContentType))
				assert.NotEmpty(t, rec.Body)
			}
		})
	}
}

func BenchmarkMergePdfs(b *testing.B) {
	h, e := setupTestHandler(nil, nil)
	body, boundary := createForm(nil, [][]byte{testdata.LoadFile(nil), testdata.LoadFile(nil)})
	req := httptest.NewRequest(echo.POST, "/pdf/merge", body)
	req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
	req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	for i := 0; i < b.N; i++ {
		_ = h.MergePdfs(ctx)
	}
}
func TestEncryptPdfFiles(t *testing.T) {
	tests := []struct {
		name               string
		files              [][]byte
		pdfTool            *pdf.PdfTool
		passwordField      FormField
		wantHttpStatusCode int
		wantContentType    string
		wantErr            bool
	}{
		{
			name:               "when just one file is provided",
			files:              [][]byte{testdata.LoadFile(t)},
			pdfTool:            nil,
			passwordField:      FormField{Field: "password", Value: "test"},
			wantHttpStatusCode: http.StatusOK,
			wantContentType:    MIMEApplicationPdf,
			wantErr:            false,
		},
		{
			name:               "when the file is invalid",
			files:              [][]byte{testdata.UnreadableFile()},
			pdfTool:            nil,
			passwordField:      FormField{Field: "password", Value: "test"},
			wantHttpStatusCode: http.StatusInternalServerError,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
		{
			name:               "when the password is not provided",
			files:              [][]byte{testdata.LoadFile(t)},
			pdfTool:            nil,
			passwordField:      FormField{Field: "password", Value: ""},
			wantHttpStatusCode: http.StatusBadRequest,
			wantContentType:    echo.MIMEApplicationJSONCharsetUTF8,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		h, e := setupTestHandler(tt.pdfTool, nil)
		body, boundary := createForm(t, tt.files, tt.passwordField)
		req := httptest.NewRequest(echo.POST, "/pdf/encrypt", body)
		req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
		req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		t.Run(tt.name, func(t *testing.T) {
			if assert.NotPanics(t, func() { _ = h.EncryptPdf(ctx) }) {
				assert.Equal(t, tt.wantHttpStatusCode, rec.Code)
				assert.Equal(t, tt.wantContentType, rec.Header().Get(echo.HeaderContentType))
				assert.NotEmpty(t, rec.Body)
			}
		})
	}
}

func BenchmarkEncryptPdf(b *testing.B) {
	h, e := setupTestHandler(nil, nil)
	body, boundary := createForm(nil, [][]byte{testdata.LoadFile(nil)}, FormField{Field: "password", Value: "test"})
	req := httptest.NewRequest(echo.POST, "/pdf/encrypt", body)
	req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
	req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	for i := 0; i < b.N; i++ {
		_ = h.EncryptPdf(ctx)
	}
}

type FormField struct {
	Field string
	Value string
}

func createForm(t *testing.T, files [][]byte, fields ...FormField) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if len(files) == 1 {
		part, _ := writer.CreateFormFile("file", "file.pdf")
		if _, err := part.Write(files[0]); err != nil {
			t.Fatal("could not write form-data")
		}
	} else {
		for i, file := range files {
			part, _ := writer.CreateFormFile("files", fmt.Sprintf("file%d.pdf", i))
			if _, err := part.Write(file); err != nil {
				t.Fatal("could not write form-data")
			}
		}
	}

	for _, field := range fields {
		err := writer.WriteField(field.Field, field.Value)
		if err != nil {
			t.Fatal("could not write form-data")
		}
	}

	return body, writer.Boundary()
}
