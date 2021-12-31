package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/headers"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type ContextKey string

const ContextKeyDevice ContextKey = "device"

// AuthMiddleware checks if a request contains the x-schism-token header and attaches the according device
type AuthMiddleware struct {
	Database *db.Sqlite
}

// NewAuthMiddleware creates a new middleware instance
func NewAuthMiddleware(database *db.Sqlite) *AuthMiddleware {
	return &AuthMiddleware{Database: database}
}

func (m *AuthMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if device is authenticated
			device, _, err := m.getAuthenticatedDevice(r)

			// Reject on any error
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Attach authenticated device to request context
			ctxWithUser := context.WithValue(r.Context(), ContextKeyDevice, device)
			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		})
	}
}

func (m *AuthMiddleware) getAuthenticatedDevice(r *http.Request) (*business.Device, int, error) {
	token := r.Header.Get(headers.HeaderSchismToken)

	// Authenticate with accesstoken
	accesstoken := &business.Accesstoken{Identifyable: db.Identifyable{Database: m.Database}}
	accesstoken, status, err := accesstoken.Authenticate(token)
	if err != nil {
		return nil, status, err
	}

	// Read device via accesstokens device_id
	device := &business.Device{Identifyable: db.Identifyable{Id: &accesstoken.DeviceId, Database: m.Database}}
	device, status, err = device.Read()
	if err != nil {
		util.Log.Panic(err.Error())
		return nil, status, err
	}

	return device, status, nil
}
