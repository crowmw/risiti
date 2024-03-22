package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/crowmw/risiti/model"
	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/component"
	"github.com/crowmw/risiti/view/uploadForm"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/validator.v2"
)

type ReceiptHandler struct {
	ReceiptService service.IReceiptService
	FileStorage    service.IFileStorage
}

func NewReceiptHandler(s service.IReceiptService, fs service.IFileStorage) *ReceiptHandler {
	return &ReceiptHandler{
		ReceiptService: s,
		FileStorage:    fs,
	}
}

func (h *ReceiptHandler) GetReceipts(w http.ResponseWriter, r *http.Request) {
	receipts, err := h.ReceiptService.ReadAll()
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
	RenderView(w, r, component.ReceiptsList(receipts), "/")
}

func (h *ReceiptHandler) SearchReceipts(w http.ResponseWriter, r *http.Request) {
	// sanitize
	s := bluemonday.UGCPolicy()
	text := s.Sanitize(r.FormValue("search"))

	receipts, err := h.ReceiptService.ReadByText(text)
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
	RenderView(w, r, component.ReceiptsList(receipts), "/")
}

func (h *ReceiptHandler) PostReceipt(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10MB

	// sanitize
	s := bluemonday.UGCPolicy()
	name, description, dateString := s.Sanitize(strings.TrimSpace(r.FormValue("name"))), s.Sanitize(strings.TrimSpace(r.FormValue("description"))), s.Sanitize(strings.TrimSpace(r.FormValue("date")))

	// parse date
	if dateString == "" {
		dateString = time.Now().Format(service.YYYYMMDD)
	}
	date, err := time.Parse(service.YYYYMMDD, dateString)
	if err != nil {
		slog.Error("Cannot parse the date", err)
		return
	}

	// File get
	file, handler, err := r.FormFile("file")
	if err != nil {
		RenderView(w, r, uploadForm.Show("File is required!"), "/upload")
		return
	}

	// create filename
	ext := path.Ext(handler.Filename)
	slug := CreateSlug(fmt.Sprintf("%s_%s", name, date.Format(service.YYYYMMDD)))
	filename := fmt.Sprintf("%s%s", slug, ext)

	// validate data
	receipt := model.Receipt{
		Name:        name,
		Description: description,
		Date:        date,
		Filename:    filename,
	}
	if err := validator.Validate(receipt); err != nil {
		RenderView(w, r, uploadForm.Show(fmt.Sprint(err)), "/upload")
		return
	}

	// name uniqnes check
	if _, err = h.ReceiptService.ReadByName(receipt.Name); err != sql.ErrNoRows {
		RenderView(w, r, uploadForm.Show("Name is already taken"), "/upload")
		return
	}

	// save file to filestorage
	err = h.FileStorage.SaveFile(file, filename)
	if err != nil {
		slog.Error("Cannot save file", err)
		return
	}

	// create receipt in db
	_, err = h.ReceiptService.Create(receipt)
	if err != nil {
		slog.Error("Cannot add receipt to storage", err)
		return
	}

	r.Header.Add("HX-Redirect", "/")
}
