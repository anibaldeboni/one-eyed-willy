package handler

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *Handler) addRoutes(e *echo.Echo) {
	e.Renderer = h.NewTemplateRegistry()
	e.GET("/", h.HomeView)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/generate", h.CreatePdfFromHtmlView)
	e.GET("/merge", h.MergePdfsView)

	e.POST("/generate", h.GeneratePdfFromHTML)
	e.POST("/merge", h.MergePdfs)
}
