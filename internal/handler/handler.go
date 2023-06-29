package handler

type Handler struct{}

const (
	MIMEApplicationPdf = "application/pdf"
)

func New() *Handler {
	return &Handler{}
}
