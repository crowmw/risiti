package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/crowmw/risiti/internal/components"
	"github.com/crowmw/risiti/internal/filestore"
	"github.com/crowmw/risiti/internal/utils"
)

const (
	YYYYMMDD = "2006-01-02"
)

type PostSubmitHandler struct {
	filestore filestore.FileStore
}

func NewPostSubmitHandler(filestore filestore.FileStore) *PostSubmitHandler {
	return &PostSubmitHandler{
		filestore: filestore,
	}
}

func (h *PostSubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Photo uploaded!")
	r.ParseMultipartForm(10 << 20) //10MB

	name := r.FormValue("name")
	dateString := r.FormValue("date")
	date, err := time.Parse(YYYYMMDD, dateString)
	if err != nil {
		slog.Error("Cannot parse the date", err)
		return
	}

	// File get
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Retrieving file!", err)
		return
	}

	// Save file to disk
	ext := path.Ext(handler.Filename)
	slug := utils.CreateSlug(fmt.Sprintf("%s_%s", name, date))
	filename := fmt.Sprintf("%s%s", slug, ext)

	err = h.filestore.SaveFile(file, filename)
	if err != nil {
		slog.Error("Saving file!", err)
		return
	}

	slog.Info("File Saved!")
	RenderView(w, r, components.Home(), "/")
}
