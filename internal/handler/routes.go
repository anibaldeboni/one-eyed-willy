package handler

import (
	"embed"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed static/*
var static embed.FS

//go:embed assets/*
var assets embed.FS

func (h *Handler) addRoutes(e *echo.Echo) {
	e.StaticFS("/static/*", echo.MustSubFS(static, "static"))
	e.StaticFS("/assets/*", echo.MustSubFS(assets, "assets"))

	e.Renderer = h.NewTemplateRegistry()

	e.GET("/", h.HomeView)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/generate", h.CreatePdfFromHtmlView)
	e.GET("/merge", h.MergePdfsView)

	e.POST("/generate", h.GeneratePdfFromHTML)
	e.POST("/merge", h.MergePdfs)
}
