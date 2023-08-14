package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/pkg/pdf"
)

type Handler struct {
	PdfRender *pdf.PdfRender
}

const (
	MIMEApplicationPdf = "application/pdf"
)

func New(e *echo.Echo) *Handler {
	h := &Handler{
		PdfRender: pdf.NewRender(),
	}
	h.addRoutes(e)
	return h
}
