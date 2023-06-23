package handler

import "github.com/labstack/echo/v4"

type createPdfFromHtmlRequest struct {
	Html string `json:"html" validate:"required,base64"`
}

func (r *createPdfFromHtmlRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}
