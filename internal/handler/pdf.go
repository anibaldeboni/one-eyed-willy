package handler

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/pkg/pdf"
	"github.com/one-eyed-willy/pkg/utils"
)

// CreatePdfFromHtml godoc
//
//	@summary		Create a pdf
//	@description	Generate a new pdf file from a html string
//	@tags			pdf
//	@accept			applcation/json
//	@produce		application/octet-stream
//	@Param			html	body	createPdfFromHTMLRequest	true	"Base64 encoded string of a html"
//	@success		200
//	@failure		400	{object}	utils.Error
//	@failure		500	{object}	utils.Error
//	@Router			/pdf [post]
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

func readBytes(files []*multipart.FileHeader) (filesBytes [][]byte, err error) {
	var wg sync.WaitGroup
	var fileBytes [][]byte

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		wg.Add(1)
		go func(src multipart.File) {
			defer wg.Done()
			buf := new(bytes.Buffer)
			if _, e := io.Copy(buf, src); err != nil {
				err = e
			}
			fileBytes = append(fileBytes, buf.Bytes())
		}(src)

	}
	wg.Wait()

	return fileBytes, err
}

// MergePdfFiles godoc
//
//	@Summary		Merge pdfs
//	@Description	Merges two or more pdfs
//	@Tags			pdf
//	@Accept			multipart/form-data
//	@Produce		application/octet-stream
//	@param			files	formData	file	true	"this is a pdf file"
//	@Success		200
//	@Failure		400	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/pdf/merge [post]
func (h *Handler) MergePdfs(c echo.Context) (err error) {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	if len(files) < 2 {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("You must provide at least two files to merge")))
	}

	fileBytes, err := readBytes(files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	pdf, err := pdf.Merge(fileBytes)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Blob(http.StatusOK, MIMEApplicationPdf, pdf)
}
