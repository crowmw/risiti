package handlers

import (
	"net/http"

	"github.com/crowmw/risiti/internal/components"
)

type GetUploadHandler struct{}

func NewGetUploadHandler() *GetUploadHandler {
	return &GetUploadHandler{}
}

func (h *GetUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.UploadForm())
}
