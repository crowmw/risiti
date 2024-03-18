package handler

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/crowmw/risiti/view/layout"

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

	err := layout.Show(layoutPath).Render(r.Context(), w)
	OnError(w, err, "Internal server error", http.StatusInternalServerError)
}

// Function to generate a slug from a string
func CreateSlug(input string) string {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace spaces with dashes
	slug := strings.ReplaceAll(processedString, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}
