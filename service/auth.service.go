package service

import (
	"errors"
	"net/http"
	"time"

	"github.com/crowmw/risiti/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
)

const (
	TOKEN_EXPIRATION_TIME = 24 * time.Hour
)

type IAuthService interface {
	GenerateToken(user model.User) (string, error)
	ValidateToken(token string) (jwt.MapClaims, error)
	GenerateCookie(user model.User) (http.Cookie, error)
}

type AuthService struct {
	JWTAuth   *jwtauth.JWTAuth
	secretKey []byte
}

func NewAuthService(secretKey []byte) *AuthService {
	jwtAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	return &AuthService{
		JWTAuth:   jwtAuth,
		secretKey: secretKey,
	}
}

func (a *AuthService) GenerateToken(user model.User) (string, error) {
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

func (a *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (a *AuthService) GenerateCookie(user model.User) (http.Cookie, error) {
	token, err := a.GenerateToken(user)
	if err != nil {
		return http.Cookie{}, err
	}

	expiration := time.Now().Add(TOKEN_EXPIRATION_TIME)
	cookie := http.Cookie{Name: "jwt", Value: token, Expires: expiration, Path: "/"}

	return cookie, nil
}
