package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api"
	"gitlab.void-ptr.org/go/schism/pkg/api/headers"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

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
			token := r.Header.Get(headers.HeaderSchismToken)
			accesstoken, device, _, err := m.getAuthenticatedDevice(r, token)

			// Reject on any error
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Attach authenticated device and token to request context
			ctxWithDeviceAndToken := context.WithValue(r.Context(), api.ContextKeyDevice, device)
			ctxWithDeviceAndToken = context.WithValue(ctxWithDeviceAndToken, api.ContextKeyToken, accesstoken)

			rWithUser := r.WithContext(ctxWithDeviceAndToken)
			next.ServeHTTP(w, rWithUser)
		})
	}
}

func (m *AuthMiddleware) getAuthenticatedDevice(r *http.Request, token string) (*business.Accesstoken, *business.Device, int, error) {

	// Authenticate with accesstoken
	accesstoken := business.NewAccesstoken(nil, m.Database)
	accesstoken, status, err := accesstoken.Authenticate(token)
	if err != nil {
		return nil, nil, status, err
	}

	// Read device via accesstokens device_id
	device := business.NewDevice(&accesstoken.DeviceId, m.Database)
	device, status, err = device.Read()
	if err != nil {
		util.Log.Panic(err.Error())
		return nil, nil, status, err
	}

	return accesstoken, device, status, nil
}
