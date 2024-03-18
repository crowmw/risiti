package handler

import (
	"net/http"

	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/home"
	"github.com/crowmw/risiti/view/signin"
	"github.com/crowmw/risiti/view/signup"
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
	RenderView(w, r, signin.Show(), "/signin")
}

func (h *UserHandler) GetSignup(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, signup.Show("", ""), "/signup")
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("confirm")

	// Check passwords is the same
	if password != passwordConfirm {
		RenderView(w, r, signup.Show(email, "Passwords are not the same."), "/signup")
		return
	}

	// Check user is already exists
	user, err := h.UserService.CheckEmail(email)
	if user || err != nil {
		RenderView(w, r, signup.Show(email, "User "+email+" already exists. Try signin."), "/signup")
		return
	}

	RenderView(w, r, home.Show(), "/")
}
