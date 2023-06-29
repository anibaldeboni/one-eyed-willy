package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/pkg/pdf"
	"github.com/one-eyed-willy/pkg/utils"
)

// CreatePdfFromHtml godoc
// @Summary Create a pdf
// @Description Generate a new pdf file from a html string
// @Tags pdf
// @Accept applcation/json
// @Produce application/octet-stream
// @Param html body createPdfFromHtmlRequest true "Base64 encoded string of a html"
// @Success 200
// @Failure 400 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /pdf [post]
func (h *Handler) GeneratePdfFromHTML(c echo.Context) error {
	req := &createPdfFromHTMLRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(req.HTML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	pdf, err := pdf.GenerateFromHTML(string(rawDecodedText))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Blob(http.StatusOK, MIMEApplicationPdf, pdf)
}
