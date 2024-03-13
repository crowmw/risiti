package handlers

import (
	"log/slog"
	"net/http"

	"github.com/crowmw/risiti/internal/components"

	"github.com/a-h/templ"
)

func onError(w http.ResponseWriter, err error, msg string, code int) {
	if err != nil {
		http.Error(w, msg, code)
		slog.Error(msg, err)
	}
}

func RenderView(w http.ResponseWriter, r *http.Request, view templ.Component) {
	if r.Header.Get("Hx-Request") == "true" {
		err := view.Render(r.Context(), w)
		onError(w, err, "Internal server error!", http.StatusInternalServerError)
		return
	}

	err := components.Layout(components.Home()).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}
