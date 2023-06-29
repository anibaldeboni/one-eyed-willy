package handler

import "github.com/labstack/echo/v4"

type createPdfFromHTMLRequest struct {
	HTML string `json:"html" validate:"required,base64"`
}

func (r *createPdfFromHTMLRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	return c.Validate(r)
}
