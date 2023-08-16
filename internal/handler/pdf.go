package handler

import (
	"bytes"
	"encoding/base64"
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
//	@Router			/pdf/generate [post]
func (h *Handler) GeneratePdfFromHTML(c echo.Context) (err error) {
	req := new(createPdfFromHTMLRequest)

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	decoded, err := base64.StdEncoding.DecodeString(req.HTML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	pdf, err := pdf.GenerateFromHTML(h.PdfRender.Context, string(decoded))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
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
//	@param			files	formData	file	true	"pdf files to merge"
//	@Success		200
//	@Failure		400	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/pdf/merge [post]
func (h *Handler) MergePdfs(c echo.Context) (err error) {
	req := new(mergePdfsRequest)

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	files, err := readBytes(req.Files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	pdf, err := pdf.Merge(files)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
}

// EncryptPdfFiles godoc
//
//		@Summary		Encrypt pdf
//		@Description	Encrypts a pdf file
//		@Tags			pdf
//		@Accept			multipart/form-data
//		@Produce		application/octet-stream
//		@param			file	formData	file	true	"pdf file to encrypt"
//	  @param			password	formData	string	true	"password to encrypt the pdf"
//		@Success		200
//		@Failure		400	{object}	utils.Error
//		@Failure		500	{object}	utils.Error
//		@Router			/pdf/encrypt [post]
func (h *Handler) EncryptPdf(c echo.Context) error {
	req := new(encryptPdfRequest)

	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}
	file, err := readBytes([]*multipart.FileHeader{req.File})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	pdf, err := pdf.Encrypt(file[0], req.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
}
