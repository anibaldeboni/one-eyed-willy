package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/internal/config"
	"github.com/one-eyed-willy/internal/web"
	"github.com/one-eyed-willy/pkg/pdf"
)

func setupTestHandler(pdfTool *pdf.PdfTool, pdfRender *pdf.PdfRender) (*Handler, *echo.Echo) {
	conf := config.InitAppConfig()
	e := web.New(conf)
	var t *pdf.PdfTool
	var r *pdf.PdfRender
	if pdfTool == nil {
		t = pdf.NewPdfTool()
	} else {
		t = pdfTool
	}
	if pdfRender == nil {
		r = pdf.NewRender()
	} else {
		r = pdfRender
	}
	h := &Handler{
		PdfRender: r,
		PdfTool:   t,
	}
	h.setupRoutes(e)

	return h, e
}
