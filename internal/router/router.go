package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/pkg/logger"
)

func New(conf *config.AppConfig) *echo.Echo {

	logger.InitWithOptions(logger.WithConfigLevel(conf.LogLevel))
	if logger.Log() != nil {
		defer logger.Log().Sync()
	}

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig(conf)))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.HideBanner = true
	e.Validator = conf.Validator

	return e
}
