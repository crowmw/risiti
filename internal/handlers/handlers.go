package handlers

import (
	"crowmw/risiti/internal/components"
	"crowmw/risiti/internal/utils"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"

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

	err := components.Layout(layoutPath).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, components.HomeView("Hello World!"), "/")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Photo uploaded!")
	r.ParseMultipartForm(10 << 20) //10MB

	name := r.FormValue("name")
	date := r.FormValue("date")

	// File get
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Retrieving file!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save file to disk
	ext := path.Ext(handler.Filename)
	slug := utils.CreateSlug(fmt.Sprintf("%s_%s", name, date))
	filename := fmt.Sprintf("%s%s", slug, ext)

	defer file.Close()
	dst, err := os.Create(fmt.Sprintf("./bin/recipes/%s", filename))
	if err != nil {
		slog.Error("Creating file!", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("File Saved!")
}
