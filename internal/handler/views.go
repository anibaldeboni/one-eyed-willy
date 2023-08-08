package handler

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed views
var viewsFS embed.FS

func (h *Handler) IndexView(c echo.Context) error {
	page, _ := viewsFS.ReadFile("views/index.html")
	return c.HTML(http.StatusOK, string(page))
}

// CreatePdfFromHtmlView godoc
func (h *Handler) CreatePdfFromHtmlView(c echo.Context) error {
	page, _ := viewsFS.ReadFile("views/generate.html")
	return c.HTML(http.StatusOK, string(page))
}

// MergePdfsView godoc
func (h *Handler) MergePdfsView(c echo.Context) error {
	page, _ := viewsFS.ReadFile("views/merge.html")
	return c.HTML(http.StatusOK, string(page))
}
