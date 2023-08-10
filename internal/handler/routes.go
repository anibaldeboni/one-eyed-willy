package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *Handler) Register(e *echo.Echo, conf *config.AppConfig) {
	e.Renderer = h.NewTemplateRegistry()
	e.GET("/", h.HomeView)
	e.GET("/docs/*", echoSwagger.WrapHandler)

	api := e.Group(conf.BaseURL)
	api.POST("", h.GeneratePdfFromHTML)
	api.GET("", h.CreatePdfFromHtmlView)
	api.POST("/merge", h.MergePdfs)
	api.GET("/merge", h.MergePdfsView)
}
