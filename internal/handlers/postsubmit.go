package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"

	"github.com/crowmw/risiti/internal/store/filestore"
	"github.com/crowmw/risiti/internal/utils"
)

type PostSubmitHandler struct {
	filestore *filestore.FileStore
}

func NewPostSubmitHandler(filestore *filestore.FileStore) *PostSubmitHandler {
	return &PostSubmitHandler{
		filestore: filestore,
	}
}

func (h *PostSubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Photo uploaded!")
	r.ParseMultipartForm(10 << 20) //10MB

	name := r.FormValue("name")
	date := r.FormValue("date")

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
}
