package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/errors"
)

// CreatePdfFromHtml godoc
//
//	@summary      Create a pdf
//	@description	Generate a new pdf file from a html string
//	@tags         PDF Generator
//	@accept       applcation/json
//	@produce      application/octet-stream
//	@Param        html	body	createPdfFromHTMLRequest	true	"Base64 encoded string of a html"
//	@success      200
//	@failure      400	{object}	errors.ValidationError
//	@Failure      413
//	@failure      500	{object}	errors.Error
//	@Router       /pdf/generate [post]
func (h *Handler) GeneratePdfFromHTML(c echo.Context) (err error) {
	req := new(createPdfFromHTMLRequest)

	if err := bindRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
	}

	decoded, err := base64.StdEncoding.DecodeString(req.HTML)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodeUndecodableBase64Err))
	}

	pdf, err := h.PdfRender.GenerateFromHTML(string(decoded))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodePdfGenerationErr))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
}

// MergePdfFiles godoc
//
//	@Summary      Merge pdfs
//	@Description  Merges two or more pdfs
//	@Tags         PDF Tools
//	@Accept       multipart/form-data
//	@Produce      application/octet-stream
//	@param        files	formData	file	true	"pdf files to merge"
//	@Success      200
//	@Failure      400	{object}	errors.ValidationError
//	@Failure      413
//	@Failure      500	{object}	errors.Error
//	@Router       /pdf/merge [post]
func (h *Handler) MergePdfs(c echo.Context) (err error) {
	req := new(mergePdfsRequest)

	if err := bindRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
	}

	files, err := readFiles(req.Files)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodeUnreadlableFileErr))
	}

	pdf, err := h.PdfTool.Merge(files)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodePdfMergeErr))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
}

// EncryptPdfFiles godoc
//
//		@Summary      Encrypt pdf
//		@Description  Encrypts a pdf file
//		@Tags         PDF Tools
//		@Accept       multipart/form-data
//		@Produce      application/octet-stream
//		@param        file		formData	file	true	"pdf file to encrypt"
//		@param        password	formData	string	true	"file password"
//		@Success      200
//		@Failure      400	{object}	errors.ValidationError
//	  @Failure      413
//		@Failure      500	{object}	errors.Error
//		@Router       /pdf/encrypt [post]
func (h *Handler) EncryptPdf(c echo.Context) error {
	req := new(encryptPdfRequest)

	if err := bindRequest(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.NewValidatorError(err))
	}
	file, err := readFile(req.File)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodeUnreadlableFileErr))
	}

	pdf, err := h.PdfTool.Encrypt(file, req.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.NewError(err, errors.CodePdfEncryptionErr))
	}

	return c.Stream(http.StatusOK, MIMEApplicationPdf, pdf)
}
