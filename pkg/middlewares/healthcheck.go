package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

var healthCheckPath = "/health"

func HealthCheck() echo.MiddlewareFunc {
	return HealthCheckWithConfig(healthCheckPath)
}

func HealthCheckWithConfig(path string) echo.MiddlewareFunc {
	healthCheckPath = path
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodGet && c.Request().URL.Path == healthCheckPath {
				return c.JSON(http.StatusOK, healthCheckResponse{Status: "OK"})
			}
			return next(c)
		}
	}
}
