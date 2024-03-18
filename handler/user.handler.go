package handler

import (
	"net/http"

	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/signin"
)

type UserHandler struct {
	UserService service.IUserService
}

func NewUserHandler(us service.IUserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

func (h *UserHandler) GetSignin(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, signin.Show(), "/")
}
