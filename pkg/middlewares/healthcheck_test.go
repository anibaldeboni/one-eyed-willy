package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	tests := []struct {
		name        string
		path        string
		response    []byte
		contentType string
	}{
		{
			name:        "when the path is /health",
			path:        healthCheckPath,
			response:    []byte{0x7b, 0x22, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3a, 0x22, 0x4f, 0x4b, 0x22, 0x7d, 0xa}, // {"status":"OK"},
			contentType: echo.MIMEApplicationJSONCharsetUTF8,
		},
		{
			name:        "when the path is other than /health",
			path:        "/other-path",
			response:    []byte("test"),
			contentType: echo.MIMETextPlainCharsetUTF8,
		},
	}
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, HealthCheck()(handler)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.contentType, rec.Header().Get(echo.HeaderContentType))
			assert.Equal(t, tt.response, rec.Body.Bytes())
		}
	}
}

func TestHealthCheckWithConfig(t *testing.T) {
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}
	tests := []struct {
		name string
		path string
	}{
		{
			name: "when the path is /health",
			path: "/health",
		},
		{
			name: "when the path is other than /health",
			path: "/my-health-check",
		},
	}
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, tt.path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		want := []byte{0x7b, 0x22, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3a, 0x22, 0x4f, 0x4b, 0x22, 0x7d, 0xa} // {"status":"OK"},

		if assert.NoError(t, HealthCheckWithConfig(tt.path)(handler)(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
			assert.Equal(t, want, rec.Body.Bytes())
		}
	}
}
