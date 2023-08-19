package handler

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/web"
	"github.com/one-eyed-willy/pkg/pdf"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	t.Run("When creating a new handler", func(t *testing.T) {
		e := echo.New()
		h := New(e)
		assert.IsType(t, &Handler{}, h)
		assert.NotNil(t, h.PdfRender)
		assert.NotNil(t, h.PdfTool)
		assert.IsType(t, &pdf.PdfRender{}, h.PdfRender)
		assert.IsType(t, &pdf.PdfTool{}, h.PdfTool)
	})

}

func TestHandlerRoutes(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "Home view route",
			method: echo.GET,
			path:   "/",
		},
		{
			name:   "Generate view route",
			method: echo.GET,
			path:   "/generate",
		},
		{
			name:   "Merge view route",
			method: echo.GET,
			path:   "/merge",
		},
		{
			name:   "Encrypt view route",
			method: echo.GET,
			path:   "/encrypt",
		},
		{
			name:   "Generate route",
			method: echo.POST,
			path:   "/pdf/generate",
		},
		{
			name:   "Merge route",
			method: echo.POST,
			path:   "/pdf/merge",
		},
		{
			name:   "Encrypt route",
			method: echo.POST,
			path:   "/pdf/encrypt",
		},
		{
			name:   "Swagger route",
			method: echo.GET,
			path:   "/docs/*",
		},
	}
	e := echo.New()
	_ = New(e)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := findRouteByPath(e.Routes(), tt.path)
			assert.NotNil(t, r)
			assert.Equal(t, tt.method, r.Method)
		})
	}
}
func findRouteByPath(routes []*echo.Route, path string) *echo.Route {
	for _, route := range routes {
		if route.Path == path {
			return route
		}
	}
	return nil
}

func setupTestHandler(pdfTool *pdf.PdfTool, pdfRender *pdf.PdfRender) (*Handler, *echo.Echo) {
	conf := config.InitAppConfig()
	e := web.New(conf)
	var t *pdf.PdfTool
	var r *pdf.PdfRender
	if pdfTool == nil {
		t = pdf.NewPdfTool()
	} else {
		t = pdfTool
	}
	if pdfRender == nil {
		r = pdf.NewRender()
	} else {
		r = pdfRender
	}
	h := &Handler{
		PdfRender: r,
		PdfTool:   t,
	}
	h.setupRoutes(e)

	return h, e
}
