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
		Path    string
		Handler func(c echo.Context) error
	}{
		{
			Path:    "/",
			Handler: h.HomeView,
		},
		{
			Path:    "/pdf",
			Handler: h.CreatePdfFromHtmlView,
		},
		{
			Path:    "/pdf/merge",
			Handler: h.MergePdfsView,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("path=%s", tt.Path), func(t *testing.T) {
			setup()
			req := httptest.NewRequest(echo.GET, tt.Path, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			if assert.NotPanics(t, func() { _ = tt.Handler(ctx) }) {
				assert.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}
