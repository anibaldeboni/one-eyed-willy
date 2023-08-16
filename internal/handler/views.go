package handler

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed views
var viewsFS embed.FS

// Define the template registry struct
type TemplateRegistry struct {
	templates Templates
}

type Templates = map[string]*template.Template

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("Template not found: %s", name)
	}

	return tmpl.ExecuteTemplate(w, name, data)
}

func (h *Handler) NewTemplateRegistry() *TemplateRegistry {
	tmpl, err := h.loadTemplates()
	if err != nil {
		panic(err)
	}
	return &TemplateRegistry{
		templates: tmpl,
	}
}

func (h *Handler) loadTemplates() (templates Templates, err error) {
	templates = make(Templates)
	files, err := viewsFS.ReadDir("views")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if (file.IsDir()) || (!file.IsDir() && (file.Name()[len(file.Name())-5:] != ".html")) {
			continue
		}

		templates[file.Name()] = template.Must(template.ParseFS(viewsFS, "views/*.html"))
	}
	return templates, nil
}

// HomeView godoc
func (h *Handler) HomeView(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", nil)
}

// CreatePdfFromHtmlView godoc
func (h *Handler) CreatePdfFromHtmlView(c echo.Context) error {
	return c.Render(http.StatusOK, "generate.html", nil)
}

// MergePdfsView godoc
func (h *Handler) MergePdfsView(c echo.Context) error {
	return c.Render(http.StatusOK, "merge.html", nil)
}

// EncryptPdfView godoc
func (h *Handler) EncryptPdfView(c echo.Context) error {
	return c.Render(http.StatusOK, "encrypt.html", nil)
}
