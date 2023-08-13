package handler

import (
	"github.com/one-eyed-willy/pkg/pdf"
)

type Handler struct {
	PdfRender *pdf.PdfRender
}

const (
	MIMEApplicationPdf = "application/pdf"
)

func New() *Handler {
	return &Handler{PdfRender: pdf.NewRender()}
}
