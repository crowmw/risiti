package service

import (
	"net/http"
	"time"

	"github.com/crowmw/risiti/model"
	"github.com/go-chi/jwtauth/v5"
)

const (
	TOKEN_EXPIRATION_TIME = 24 * time.Hour
)

type AuthService struct {
	JWTAuth   *jwtauth.JWTAuth
	secretKey []byte
}

func NewAuthService(secretKey []byte) AuthService {
	jwtAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	return AuthService{
		JWTAuth:   jwtAuth,
		secretKey: secretKey,
	}
}

func (a *AuthService) generateToken(user model.User) (string, error) {
	payload := map[string]interface{}{
		"email": user.Email,
		"id":    user.ID,
		"exp":   time.Now().Add(TOKEN_EXPIRATION_TIME).Unix(), // Token expires in 24 hours
	}

	_, tokenString, err := a.JWTAuth.Encode(payload)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) generateCookie(user model.User) (http.Cookie, error) {
	token, err := a.generateToken(user)
	if err != nil {
		return http.Cookie{}, err
	}

	expiration := time.Now().Add(TOKEN_EXPIRATION_TIME)
	cookie := http.Cookie{Name: "jwt", Value: token, Expires: expiration, Path: "/"}

	return cookie, nil
}

func (a *AuthService) SignIn(w *http.ResponseWriter, user *model.User) error {
	cookie, err := a.generateCookie(*user)
	if err != nil {
		return err
	}
	http.SetCookie(*w, &cookie)
	return nil
}

func (a *AuthService) SignOut(w *http.ResponseWriter) {
	c := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(*w, &c)
}

func (a *AuthService) Authorize(r *http.Request) error {
	_, err := jwtauth.VerifyRequest(a.JWTAuth, r, jwtauth.TokenFromCookie)
	if err != nil {
		return err
	}
	return nil
}
