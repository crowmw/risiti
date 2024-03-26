package handler

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/crowmw/risiti/model"
	"github.com/crowmw/risiti/service"
	"github.com/crowmw/risiti/view/signin"
	"github.com/crowmw/risiti/view/signup"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

type UserHandler struct {
	UserService service.IUserService
	AuthService service.AuthService
}

func NewUserHandler(us service.IUserService, jwt service.AuthService) *UserHandler {
	return &UserHandler{
		UserService: us,
		AuthService: jwt,
	}
}

func (h *UserHandler) GetSignin(w http.ResponseWriter, r *http.Request) {
	RenderView(w, r, signin.Show("", ""), "/signin")
}

func (h *UserHandler) GetSignup(w http.ResponseWriter, r *http.Request) {
	// Check is any user exists in system
	anyUserExists, err := h.UserService.AnyExists()
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !anyUserExists {
		RenderView(w, r, signup.Show("", ""), "/signup")
		return
	}

	if err := h.AuthService.Authorize(r); err != nil {
		RenderView(w, r, signin.Show("", ""), "/signin")
		return
	}

	RenderView(w, r, signup.Show("", ""), "/signup")
}

func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	// sanitize
	s := bluemonday.UGCPolicy()
	user.Email, user.Password = s.Sanitize(strings.TrimSpace(r.FormValue("email"))), s.Sanitize(strings.TrimSpace(r.FormValue("password")))

	// validate
	if err := validator.Validate(user); err != nil {
		RenderView(w, r, signup.Show(user.Email, fmt.Sprint(err)), "/signup")
		return
	}

	passwordConfirm := r.FormValue("confirm")

	// Check passwords is the same
	if user.Password != passwordConfirm {
		RenderView(w, r, signup.Show(user.Email, "Passwords are not the same."), "/signup")
		return
	}

	// Check user is already exists
	if _, err := h.UserService.Read(user.Email); err != sql.ErrNoRows {
		RenderView(w, r, signup.Show(user.Email, "User "+user.Email+" already exists. Try signin."), "/signup")
		return
	}

	// Create new user
	newUser, err := h.UserService.Create(user)
	if err != nil {
		RenderView(w, r, signup.Show(user.Email, "Something went wrong while creating user."), "/signup")
		slog.Error(fmt.Sprint(err))
		return
	}

	// Signin new user
	err = h.AuthService.SignIn(&w, &newUser)
	if err != nil {
		RenderView(w, r, signup.Show(user.Email, "Something went wrong while signing in new user."), "/signup")
		return
	}

	w.Header().Add("HX-Redirect", "/")
}

func (h *UserHandler) PostSignin(w http.ResponseWriter, r *http.Request) {
	// sanitize
	s := bluemonday.UGCPolicy()
	email := s.Sanitize(strings.TrimSpace(r.FormValue("email")))
	password := strings.TrimSpace(r.FormValue("password"))

	// Check user is already exists
	existedUser, err := h.UserService.Read(email)
	if err != nil {
		RenderView(w, r, signin.Show(email, "Wrong username or password."), "/signin")
		return
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(existedUser.Password),
		[]byte(password),
	); err != nil {
		RenderView(w, r, signin.Show(email, "Wrong username or password."), "/signin")
		return
	}

	err = h.AuthService.SignIn(&w, &existedUser)
	if err != nil {
		RenderView(w, r, signin.Show(email, "Something went wrong while signing in user."), "/signin")
		return
	}

	w.Header().Add("HX-Redirect", "/")
}

func (h *UserHandler) GetSignout(w http.ResponseWriter, r *http.Request) {
	h.AuthService.SignOut(&w)
	w.Header().Add("HX-Redirect", "/")
}
