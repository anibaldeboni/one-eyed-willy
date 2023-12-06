package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var logger = zap.NewNop()

func TestHealthCheck(t *testing.T) {
	t.Run("health check", func(t *testing.T) {
		e := New(logger)
		req := httptest.NewRequest(echo.GET, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, "{\"status\":\"OK\"}\n", rec.Body.String())
	})
}
