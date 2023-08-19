package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/pkg/middlewares"
	"go.uber.org/zap"
)

func New(logger *zap.SugaredLogger) *echo.Echo {

	if logger != nil {
		//nolint:errcheck
		defer logger.Sync()
	}

	e := echo.New()
	e.Binder = NewFileBinder(e.Binder)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Infof("method=%s path=%s status=%d request_id=%s", v.Method, v.URI, v.Status, v.RequestID)
			} else {
				logger.Errorf("method=%s path=%s status=%d request_id=%s error=%v", v.Method, v.URI, v.Status, v.RequestID, v.Error)
			}
			return nil
		},
	}))
	// e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig()))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middlewares.HealthCheck())
	e.HideBanner = true
	e.Validator = NewValidator()

	return e
}
