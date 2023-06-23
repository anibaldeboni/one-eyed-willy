package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/router"
)

var (
	h *Handler
	e *echo.Echo
)

func setup() {
	conf, _ := config.InitAppConfig()
	h = &Handler{}
	e = router.New(conf)
}
