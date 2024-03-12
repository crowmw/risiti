package handlers

import (
	"crowmw/risiti/internal/components"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func onError(w http.ResponseWriter, err error, msg string, code int) {
	if err != nil {
		http.Error(w, msg, code)
		log.Println(msg, err)
	}
}

func RenderView(w http.ResponseWriter, r *http.Request, view templ.Component, layoutPath string) {
	if r.Header.Get("Hx-Request") == "true" {
		err := view.Render(r.Context(), w)
		onError(w, err, "Internal server error!", http.StatusInternalServerError)
		return
	}

	err := components.Layout(layoutPath, "home").Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.HomeView("Hello World!"), "/")
}
