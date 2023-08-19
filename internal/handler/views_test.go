package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestViews(t *testing.T) {
	tests := []struct {
		Path string
	}{
		{
			Path: "/",
		},
		{
			Path: "/generate",
		},
		{
			Path: "/merge",
		},
		{
			Path: "/encrypt",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("path=%s", tt.Path), func(t *testing.T) {
			_, e := setupTestHandler(nil, nil)
			req := httptest.NewRequest(echo.GET, tt.Path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMETextHTMLCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
			assert.NotEmpty(t, rec.Body)
		})
	}
}

func TestTemplateRegistry(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "when template is valid",
			path:    "home.html",
			wantErr: false,
		},
		{
			name:    "when template is invalid",
			path:    "madonna-standing-in-the-front-yard.html",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, e := setupTestHandler(nil, nil)
			tr := h.NewTemplateRegistry()
			buf := new(bytes.Buffer)
			err := tr.Render(buf, tt.path, nil, e.NewContext(nil, nil))

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantErr, buf.Len() == 0)
		})
	}
}
