package handler

import "github.com/labstack/echo/v4"

func (h *Handler) Register(api *echo.Group) {
	api.POST("", h.GeneratePdfFromHTML)
	api.POST("/merge", h.MergePdfs)
}
