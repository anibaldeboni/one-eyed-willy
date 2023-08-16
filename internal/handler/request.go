package handler

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

// swagger:parameters createPdfFromHTMLRequest
type createPdfFromHTMLRequest struct {
	HTML string `json:"html" validate:"required,base64"`
}

func (r *createPdfFromHTMLRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.Validate(r)
}

type mergePdfsRequest struct {
	Files []*multipart.FileHeader `form:"files" validate:"required,gt=1"`
}

func (r *mergePdfsRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.Validate(r)
}

type encryptPdfRequest struct {
	File     *multipart.FileHeader `form:"file" validate:"required,lt=2"`
	Password string                `form:"password" validate:"required,gt=0"`
}

func (r *encryptPdfRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.Validate(r)
}
