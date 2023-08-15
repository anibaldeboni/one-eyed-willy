package handler

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed static/*
var static embed.FS

func (h *Handler) addRoutes(e *echo.Echo) {
	staticHandler := echo.WrapHandler(http.FileServer(http.FS(static)))
	// The embedded files will all be in the '/static' folder so need to rewrite the request (could also do this with fs.Sub)
	staticRewrite := middleware.Rewrite(map[string]string{"/*": "/static/$1"})

	e.Renderer = h.NewTemplateRegistry()

	e.GET("/*", staticHandler, staticRewrite)
	e.GET("/", h.HomeView)
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.GET("/generate", h.CreatePdfFromHtmlView)
	e.GET("/merge", h.MergePdfsView)

	e.POST("/generate", h.GeneratePdfFromHTML)
	e.POST("/merge", h.MergePdfs)
}
