package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/web"
)

var (
	h *Handler
	e *echo.Echo
)

func setup() {
	conf := config.InitAppConfig()
	h = New()
	e = web.New(conf)
}
