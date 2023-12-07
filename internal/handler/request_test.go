package handler

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/testdata"
)

func TestCreatePdfFromHTMLRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		header  string
		footer  string
		wantErr bool
	}{
		{
			name:    "When all fields are  valid",
			body:    "aGVsbG8gd29ybGQ=",
			header:  "aGVsbG8gd29ybGQ=",
			footer:  "aGVsbG8gd29ybGQ=",
			wantErr: false,
		},
		{
			name:    "When body is invalid",
			body:    "invalid-html",
			header:  "aGVsbG8gd29ybGQ=",
			footer:  "aGVsbG8gd29ybGQ=",
			wantErr: true,
		},
		{
			name:    "When header is invalid",
			body:    "aGVsbG8gd29ybGQ=",
			header:  "invalid-html",
			footer:  "aGVsbG8gd29ybGQ=",
			wantErr: true,
		},
		{
			name:    "When footer is invalid",
			body:    "aGVsbG8gd29ybGQ=",
			header:  "aGVsbG8gd29ybGQ=",
			footer:  "invalid-html",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			req := httptest.NewRequest(echo.POST, "/", strings.NewReader(fmt.Sprintf(`{"html":"%s", "headerTemplate": "%s", "footerTemplate": "%s"}`, tt.body, tt.header, tt.footer)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)
			if err := bindRequest(ctx, new(createPdfFromHTMLRequest)); (err != nil) != tt.wantErr {
				t.Errorf("createPdfFromHTMLRequestBind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMergePdfRequest(t *testing.T) {
	tests := []struct {
		name    string
		files   [][]byte
		wantErr bool
	}{
		{
			name:    "When request is valid",
			files:   [][]byte{testdata.LoadFile(t), testdata.LoadFile(t)},
			wantErr: false,
		},
		{
			name:    "When request is invalid",
			files:   [][]byte{testdata.LoadFile(t)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			body, boundary := createForm(t, tt.files)
			req := httptest.NewRequest(echo.POST, "/", body)
			req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
			req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)
			if err := bindRequest(ctx, new(mergePdfsRequest)); (err != nil) != tt.wantErr {
				t.Errorf("mergePdfRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptPdfRequest(t *testing.T) {
	type args struct {
		File     [][]byte
		Password string
	}
	tests := []struct {
		name    string
		request args
		wantErr bool
	}{
		{
			name:    "When request is valid",
			request: args{File: [][]byte{testdata.LoadFile(t)}, Password: "password"},
			wantErr: false,
		},
		{
			name:    "When password is empty",
			request: args{File: [][]byte{testdata.LoadFile(t)}, Password: ""},
			wantErr: true,
		},
		{
			name:    "When file is empty",
			request: args{File: [][]byte{}, Password: "password"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			body, boundary := createForm(t, tt.request.File, FormField{Field: "password", Value: tt.request.Password})
			req := httptest.NewRequest(echo.POST, "/", body)
			req.Header.Set(echo.HeaderContentType, fmt.Sprintf("%s; boundary=%s", echo.MIMEMultipartForm, boundary))
			req.Header.Set(echo.HeaderContentLength, fmt.Sprintf("%d", len(body.Bytes())))
			res := httptest.NewRecorder()
			ctx := e.NewContext(req, res)

			if err := bindRequest(ctx, new(encryptPdfRequest)); (err != nil) != tt.wantErr {
				t.Errorf("encryptPdfRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
