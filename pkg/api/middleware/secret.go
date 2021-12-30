package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
)

// SecretMiddleware checks if the x-schism-secret header is containing the API secret
type SecretMiddleware struct {
	ApiSecret string
}

// NewSecretMiddleware creates a new middleware instance
func NewSecretMiddleware(secret string) *SecretMiddleware {
	return &SecretMiddleware{secret}
}

func (m *SecretMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := r.Header.Get("x-schism-secret")
			// Reject if secrets do not match
			if secret != m.ApiSecret {
				http.Error(w, errors.StatusUnauthorized, http.StatusUnauthorized)
				return
			}
			// Check passed, secret is fine
			next.ServeHTTP(w, r)
		})
	}
}
