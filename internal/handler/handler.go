package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/one-eyed-willy/pkg/pdf"
	"go.uber.org/zap"
)

type Handler struct {
	PdfRender *pdf.PdfRender
	PdfTool   *pdf.PdfTool
}

const (
	MIMEApplicationPdf = "application/pdf"
)

func New(e *echo.Echo, logger *zap.Logger) *Handler {
	h := &Handler{
		PdfRender: pdf.NewRender(logger),
		PdfTool:   pdf.NewPdfTool(),
	}
	h.setupRoutes(e)
	return h
}
