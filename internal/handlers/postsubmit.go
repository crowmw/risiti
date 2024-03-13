package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/crowmw/risiti/internal/components"
	"github.com/crowmw/risiti/internal/filestore"
	receiptrepo "github.com/crowmw/risiti/internal/repo"
	"github.com/crowmw/risiti/internal/tools"
)

const (
	YYYYMMDD = "2006-01-02"
)

type PostSubmitHandler struct {
	filestore   *filestore.FileStore
	receiptrepo receiptrepo.IReceiptRepo
}

func NewPostSubmitHandler(filestore *filestore.FileStore, receiptrepo receiptrepo.IReceiptRepo) *PostSubmitHandler {
	return &PostSubmitHandler{
		filestore:   filestore,
		receiptrepo: receiptrepo,
	}
}

func (h *PostSubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10MB

	name := r.FormValue("name")
	dateString := r.FormValue("date")
	description := r.FormValue("description")
	date, err := time.Parse(YYYYMMDD, dateString)
	if err != nil {
		slog.Error("Cannot parse the date", err)
		return
	}

	// File get
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Cannot retrieve file from formdata", err)
		return
	}

	// Save file to disk
	ext := path.Ext(handler.Filename)
	slug := tools.CreateSlug(fmt.Sprintf("%s_%s", name, date.Format(YYYYMMDD)))
	filename := fmt.Sprintf("%s%s", slug, ext)

	err = h.filestore.SaveFile(file, filename)
	if err != nil {
		slog.Error("Cannot save file", err)
		return
	}

	receipt := receiptrepo.Receipt{
		Name:        name,
		Description: description,
		Date:        date,
		Filename:    filename,
	}
	err = h.receiptrepo.Add(receipt)
	if err != nil {
		slog.Error("Cannot add receipt to storage", err)
		return
	}

	RenderView(w, r, components.Home(), "/")
}
