package handlers

import (
	"log/slog"
	"net/http"

	"github.com/crowmw/risiti/internal/components"

	"github.com/a-h/templ"
)

func OnError(w http.ResponseWriter, err error, msg string, code int) {
	if err != nil {
		http.Error(w, msg, code)
		slog.Error(msg, err)
	}
}

func RenderView(w http.ResponseWriter, r *http.Request, view templ.Component, layoutPath string) {
	if r.Header.Get("Hx-Request") == "true" {
		err := view.Render(r.Context(), w)
		OnError(w, err, "Internal server error!", http.StatusInternalServerError)
		return
	}

	err := components.Layout(layoutPath).Render(r.Context(), w)
	OnError(w, err, "Internal server error", http.StatusInternalServerError)
}
