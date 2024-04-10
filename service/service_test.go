package service

import (
	"net/http"
	"testing"

	"github.com/crowmw/risiti/mocks"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, &AuthTestSuite{})
}

func (ats *AuthTestSuite) TestSignIn() {
	authService := NewAuthService([]byte("secretKey"))
	
	authService.SignIn()

}