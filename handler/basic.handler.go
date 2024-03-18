package handler

import (
	"net/http"

	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/home"
	"github.com/crowmw/risiti/view/uploadForm"
)

type BasicHandler struct {
	ReceiptService service.IReceiptService
}

func NewBasicHandler(s service.IReceiptService) *BasicHandler {
	return &BasicHandler{
		ReceiptService: s,
	}
}

func (h *BasicHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, home.Show(), "/")
}

func (h *BasicHandler) GetUpload(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, uploadForm.Show(""), "/upload")
}
