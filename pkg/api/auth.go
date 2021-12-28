package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AuthMiddleware struct {
	token *string
}

func NewAuthMiddleware(token *string) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (m *AuthMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}
