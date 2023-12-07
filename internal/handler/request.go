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
	HTML            string  `json:"html" validate:"required,base64,gt=0"`
	Landscape       bool    `json:"landscape"`
	PrintBackground bool    `json:"printBackground"`
	OmitBackground  bool    `json:"omitBackground"`
	Scale           float64 `json:"scale"`
	PaperWidth      float64 `json:"paperWidth"`
	PaperHeight     float64 `json:"paperHeight"`
	MarginTop       float64 `json:"marginTop"`
	MarginBottom    float64 `json:"marginBottom"`
	MarginLeft      float64 `json:"marginLeft"`
	MarginRight     float64 `json:"marginRight"`
	HeaderTemplate  string  `json:"headerTemplate" validate:"omitempty,base64,gt=0"`
	FooterTemplate  string  `json:"footerTemplate" validate:"omitempty,base64,gt=0"`
}

type mergePdfsRequest struct {
	Files []*multipart.FileHeader `form:"files" validate:"required,gt=1"`
}

type encryptPdfRequest struct {
	File     *multipart.FileHeader `form:"file" validate:"required"`
	Password string                `form:"password" validate:"required,gt=0"`
}
