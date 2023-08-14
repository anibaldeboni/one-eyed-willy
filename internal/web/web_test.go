package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	t.Run("health check", func(t *testing.T) {
		_, rec := setupServer(
			httptest.NewRequest(echo.GET, "/health", nil),
		)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupServer(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	conf := config.InitAppConfig()
	e := New(conf)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
