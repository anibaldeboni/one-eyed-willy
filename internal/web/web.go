package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/pkg/logger"
	"github.com/one-eyed-willy/pkg/middlewares"
)

func New(conf *config.AppConfig) *echo.Echo {

	logger.InitWithOptions(logger.WithConfigLevel(conf.LogLevel))
	if logger.Log() != nil {
		//nolint:errcheck
		defer logger.Log().Sync()
	}

	e := echo.New()
	e.Binder = NewFileBinder(e.Binder)
	e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig()))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middlewares.HealthCheck())
	e.HideBanner = true
	e.Validator = conf.Validator

	return e
}
