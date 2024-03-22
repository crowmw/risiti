package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/crowmw/risiti/model"
	"github.com/crowmw/risiti/sanitize"
	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/component"
	"github.com/crowmw/risiti/view/home"
	"github.com/crowmw/risiti/view/uploadForm"
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
	text, err := sanitize.SanitizeHTML(r.FormValue("search"))
	if err != nil {
		OnError(w, err, "Search text is invalid!", http.StatusBadRequest)
		return
	}

	receipts, err := h.ReceiptService.ReadByText(text)
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
	RenderView(w, r, component.ReceiptsList(receipts), "/")
}

func (h *ReceiptHandler) PostReceipt(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) //10MB

	name, err := sanitize.SanitizeHTML(r.FormValue("name"))
	if name == "" {
		RenderView(w, r, uploadForm.Show("Name cannot be empty!"), "/upload")
		return
	}
	if err != nil {
		OnError(w, err, "Name is invalid!", http.StatusBadRequest)
		return
	}

	description, err := sanitize.SanitizeHTML(r.FormValue("description"))
	if err != nil {
		OnError(w, err, "Description is invalid!", http.StatusBadRequest)
		return
	}

	dateString, err := sanitize.SanitizeHTML(r.FormValue("date"))
	if err != nil {
		OnError(w, err, "Date is invalid!", http.StatusBadRequest)
		return
	}

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

	// Save file to disk
	ext := path.Ext(handler.Filename)
	slug := CreateSlug(fmt.Sprintf("%s_%s", name, date.Format(service.YYYYMMDD)))
	filename := fmt.Sprintf("%s%s", slug, ext)

	err = h.FileStorage.SaveFile(file, filename)
	if err != nil {
		slog.Error("Cannot save file", err)
		return
	}

	receipt := model.Receipt{
		Name:        name,
		Description: description,
		Date:        date,
		Filename:    filename,
	}
	_, err = h.ReceiptService.Create(receipt)
	if err != nil {
		slog.Error("Cannot add receipt to storage", err)
		return
	}

	RenderView(w, r, home.Show(), "/")
}
