package handler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
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
func (h *Handler) GeneratePdfFromHTML(c echo.Context) (err error) {
	req := &createPdfFromHTMLRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(req.HTML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	ch := make(chan []byte)

	go func() {
		defer close(ch)
		pdf, err := pdf.GenerateFromHTML(string(rawDecodedText))
		if err == nil {
			ch <- pdf
		}
	}()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Blob(http.StatusOK, MIMEApplicationPdf, <-ch)
}

func (h *Handler) MergePdfs(c echo.Context) (err error) {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	if len(files) < 2 {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("You must provide at least two files to merge")))
	}

	var fileBytes [][]byte
	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, src); err != nil {
			return err
		}

		fileBytes = append(fileBytes, buf.Bytes())
	}

	ch := make(chan []byte)

	go func() {
		defer close(ch)
		pdf, err := pdf.Merge(fileBytes)
		if err == nil {
			ch <- pdf
		}
	}()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Blob(http.StatusOK, MIMEApplicationPdf, <-ch)
}
