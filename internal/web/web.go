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
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middlewares.HealthCheck())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		LogMethod:    true,
		LogRequestID: true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Infof("method=%s path=%s status=%d request_id=%s latency=%d", v.Method, v.URI, v.Status, v.RequestID, v.Latency)
			} else {
				logger.Errorf("method=%s path=%s status=%d request_id=%s latency=%d error=%v", v.Method, v.URI, v.Status, v.RequestID, v.Latency, v.Error)
			}
			return nil
		},
	}))
	e.HideBanner = true
	e.Validator = NewValidator()

	return e
}
