package handler

import (
	"net/http"

	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/home"
	"github.com/crowmw/risiti/view/signin"
	"github.com/crowmw/risiti/view/signup"
	"github.com/crowmw/risiti/view/uploadForm"
)

type BasicHandler struct {
	ReceiptService service.ReceiptService
	UserService    service.UserService
	AuthService    service.AuthService
}

func NewBasicHandler(s service.ReceiptService, u service.UserService, a service.AuthService) *BasicHandler {
	return &BasicHandler{
		ReceiptService: s,
		UserService:    u,
		AuthService:    a,
	}
}

func (h *BasicHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	// Check is any user exists in system
	anyUserExists, err := h.UserService.GetAny()
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !anyUserExists {
		RenderView(w, r, signup.Show("", ""), "/signup")
		return
	}

	// Authorize request
	if err := h.AuthService.Authorize(r); err != nil {
		RenderView(w, r, signin.Show("", ""), "/signin")
		return
	}

	RenderView(w, r, home.Show(), "/")
}

func (h *BasicHandler) GetUpload(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, uploadForm.Show(""), "/upload")
}
