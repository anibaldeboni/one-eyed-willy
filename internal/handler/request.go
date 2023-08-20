package handler

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func bindRequest(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	return c.Validate(req)
}

// swagger:parameters createPdfFromHTMLRequest
type createPdfFromHTMLRequest struct {
	HTML string `json:"html" validate:"required,base64,gt=0"`
}

type mergePdfsRequest struct {
	Files []*multipart.FileHeader `form:"files" validate:"required,gt=1"`
}

type encryptPdfRequest struct {
	File     *multipart.FileHeader `form:"file" validate:"required,lt=2"`
	Password string                `form:"password" validate:"required,gt=0"`
}
