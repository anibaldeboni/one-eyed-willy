package handler

import (
	"context"

	"github.com/one-eyed-willy/pkg/pdf"
)

type Handler struct {
	Context context.Context
	Cancel  context.CancelFunc
}

const (
	MIMEApplicationPdf = "application/pdf"
)

func New() (*Handler, error) {
	ctx, cancel, err := pdf.NewContext()
	if err != nil {
		return nil, err
	}
	return &Handler{Context: ctx, Cancel: cancel}, nil
}
