package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePdfCaseSuccess(t *testing.T) {
	setup()
	var (
		reqJSON = `{"html":"PGh0bWw+CjxoZWFkPgoJPHRpdGxlPk15IFBERiBGaWxlPC90aXRsZT4KPC9oZWFkPgo8Ym9keT4KCTxwPkhlbGxvIHRoZXJlISBJJ20gYSBwZGYgZmlsZSBnZW5lcmF0ZSBmcm9tIGEgaHRtbCB1c2luZyBnbyBhbmQgZ29wZGYgcGFja2FnZTwvcD4KPC9ib2R5Pgo8L2h0bWw+"}`
	)
	req := httptest.NewRequest(echo.POST, "/pdf", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GeneratePdfFromHtml(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		assert.Equal(t, "application/pdf", rec.Header().Clone().Get("Content-Type"))
		assert.Greater(t, len(rec.Body.Bytes()), 0)
		assert.NotEmpty(t, rec.Body)
	}
}
