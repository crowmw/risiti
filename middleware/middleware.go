package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/crowmw/risiti/handler"
	"github.com/crowmw/risiti/view/signin"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func generateRandomString(length int) string {

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	htmxNonce := generateRandomString(16)
	responseTargetsNonse := generateRandomString(16)
	twNonce := generateRandomString(16)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set then in context
		ctx := context.WithValue(r.Context(), "htmxNonce", htmxNonce)
		ctx = context.WithValue(ctx, "twNonce", twNonce)
		ctx = context.WithValue(ctx, "responseTargetsNonse", responseTargetsNonse)

		// the hash of the CSS that HTMX injects
		htmxCSSHash := "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg="

		cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' 'nonce-%s' ; style-src 'nonce-%s' '%s';", htmxNonce, responseTargetsNonse, twNonce, htmxCSSHash)
		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var CORS = cors.Handler(cors.Options{
	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	AllowedOrigins: []string{"https://*", "http://*"},
	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: false,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
})

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())

			if err != nil {
				handler.RenderView(w, r, signin.Show("", ""), "/signin")
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				handler.RenderView(w, r, signin.Show(fmt.Sprint(claims["email"]), ""), "/signin")
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
