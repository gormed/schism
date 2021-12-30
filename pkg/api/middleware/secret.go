package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
)

type SecretMiddleware struct {
	ApiSecret string
}

func NewSecretMiddleware(secret string) *SecretMiddleware {
	return &SecretMiddleware{secret}
}

func (m *SecretMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := r.Header.Get("x-schism-secret")
			if secret != m.ApiSecret {
				http.Error(w, errors.StatusUnauthorized, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
