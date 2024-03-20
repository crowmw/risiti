package handler

import (
	"net/http"

	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/home"
	"github.com/crowmw/risiti/view/uploadForm"
)

type BasicHandler struct {
	ReceiptService service.IReceiptService
	UserService    service.IUserService
}

func NewBasicHandler(s service.IReceiptService, u service.IUserService) *BasicHandler {
	return &BasicHandler{
		ReceiptService: s,
		UserService:    u,
	}
}

func (h *BasicHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	anyUserExists, err := h.UserService.AnyExists()
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if anyUserExists {
		RenderView(w, r, home.Show(), "/")
		return
	}

}

func (h *BasicHandler) GetUpload(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, uploadForm.Show(""), "/upload")
}
