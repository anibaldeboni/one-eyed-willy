package handler

import (
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
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("path=%s", tt.Path), func(t *testing.T) {
			setup()
			req := httptest.NewRequest(echo.GET, tt.Path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMETextHTMLCharsetUTF8, rec.Header().Clone().Get("Content-Type"))
		})
	}
}
