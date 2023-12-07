package web

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/one-eyed-willy/pkg/middlewares"
	"go.uber.org/zap"
)

func New(logger *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Binder = NewFileBinder(e.Binder)
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middlewares.HealthCheck())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogError:      true,
		LogLatency:    true,
		LogMethod:     true,
		LogRequestID:  true,
		HandleError:   true,
		LogValuesFunc: logValues(logger),
	}))
	e.HideBanner = true
	e.Validator = NewValidator()

	return e
}

func logValues(logger *zap.Logger) func(c echo.Context, v middleware.RequestLoggerValues) error {
	return func(c echo.Context, v middleware.RequestLoggerValues) error {
		method := zap.String("method", v.Method)
		path := zap.String("path", v.URI)
		status := zap.Int("status", v.Status)
		requestID := zap.String("request_id", v.RequestID)
		latency := zap.Int64("latency", int64(v.Latency/time.Millisecond))

		if v.Error != nil {
			logger.Error("REQUEST", method, path, status, requestID, latency, zap.Error(v.Error))
		} else {
			logger.Info("REQUEST", method, path, status, requestID, latency)
		}

		return nil
	}
}
