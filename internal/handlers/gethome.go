package handlers

import (
	"net/http"

	"github.com/crowmw/risiti/internal/components"
)

type GetHomeHandler struct {
}

func NewGetHomeHandler() *GetHomeHandler {
	return &GetHomeHandler{}
}

func (h *GetHomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.Home(), "/")
}
