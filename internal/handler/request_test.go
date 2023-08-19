package handler

import (
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/testdata"
)

func TestCreatePdfFromHTMLRequestBind(t *testing.T) {
	tests := []struct {
		name    string
		request createPdfFromHTMLRequest
		wantErr bool
	}{
		{
			name:    "When request is valid",
			request: createPdfFromHTMLRequest{HTML: "aGVsbG8gd29ybGQ="},
			wantErr: false,
		},
		{
			name:    "When request is invalid",
			request: createPdfFromHTMLRequest{HTML: "invalid-html"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			req := httptest.NewRequest(echo.POST, "/", strings.NewReader(fmt.Sprintf(`{"html":"%s"}`, tt.request.HTML)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)
			if err := tt.request.bind(ctx); (err != nil) != tt.wantErr {
				t.Errorf("createPdfFromHTMLRequestBind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptPdfRequest(t *testing.T) {
	tests := []struct {
		name    string
		request encryptPdfRequest
		wantErr bool
	}{
		{
			name:    "When request is valid",
			request: encryptPdfRequest{File: &multipart.FileHeader{}, Password: "password"},
			wantErr: false,
		},
		{
			name:    "When request is invalid",
			request: encryptPdfRequest{File: &multipart.FileHeader{}, Password: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			body, boundary := createForm(t, [][]byte{testdata.LoadFile(t)})
			req := httptest.NewRequest(echo.POST, "/", body)
			req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
			req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)
			if err := tt.request.bind(ctx); (err != nil) != tt.wantErr {
				t.Errorf("encryptPdfRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
