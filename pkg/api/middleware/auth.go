package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
)

type ContextKey string

const ContextKeyDevice ContextKey = "device"

// AuthMiddleware checks if a request contains the x-schism-token header and attaches the according device
type AuthMiddleware struct {
}

// NewAuthMiddleware creates a new middleware instance
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if device is authenticated
			device, err := GetAuthenticatedDevice(r)
			// Reject on any error
			if err != nil {
				http.Error(w, errors.StatusUnauthorized, http.StatusUnauthorized)
				return
			}
			// Attach authenticated device to request context
			ctxWithUser := context.WithValue(r.Context(), ContextKeyDevice, device)
			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		})
	}
}

func GetAuthenticatedDevice(r *http.Request) (*business.Device, error) {
	token := r.Header.Get("x-schism-token")
	// Authenticate with accesstoken
	accesstoken := &business.Accesstoken{}
	accesstoken, _, err := accesstoken.Authenticate(token)
	if err != nil {
		return nil, err
	}
	// Read device via accesstokens device_id
	device := &business.Device{Identifyable: db.Identifyable{Id: &accesstoken.DeviceId}}
	device, _, err = device.Read()
	if err != nil {
		return nil, err
	}
	return device, nil
}
