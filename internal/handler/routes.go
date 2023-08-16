package handler

import (
	"embed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed static/*
var static embed.FS

//go:embed assets/*
var assets embed.FS

// maxFileSize is the maximum size of a file that can be uploaded
const maxFileSize = "15M"

func (h *Handler) setupRoutes(e *echo.Echo) {
	e.StaticFS("/static/*", echo.MustSubFS(static, "static"))
	e.StaticFS("/assets/*", echo.MustSubFS(assets, "assets"))

	e.Renderer = h.NewTemplateRegistry()

	e.GET("/", h.HomeView)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/generate", h.CreatePdfFromHtmlView)
	e.GET("/merge", h.MergePdfsView)

	api := e.Group("/pdf")
	api.Use(middleware.BodyLimit(maxFileSize))
	api.POST("/generate", h.GeneratePdfFromHTML)
	api.POST("/merge", h.MergePdfs)
}
