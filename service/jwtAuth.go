package service

import (
	"github.com/crowmw/risiti/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
)

type IJWTAuth interface {
	GenerateToken(user model.User) (string, error)
	ValidateToken(token string) (jwt.MapClaims, error)
}

type JWTAuth struct {
	JWTAuth   *jwtauth.JWTAuth
	secretKey []byte
}

func NewJWTAuth(secretKey []byte) *JWTAuth {
	jwtAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	return &JWTAuth{
		JWTAuth:   jwtAuth,
		secretKey: secretKey,
	}
}

func (a *JWTAuth) GenerateToken(user model.User) (string, error) {
	return "", nil
}

func (a *JWTAuth) ValidateToken(token string) (jwt.MapClaims, error) {
	return jwt.MapClaims{}, nil
}
