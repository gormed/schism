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
			device, err := GetAuthenticatedDevice(r)
			if err != nil {
				http.Error(w, errors.StatusUnauthorized, http.StatusUnauthorized)
				return
			}

			ctxWithUser := context.WithValue(r.Context(), ContextKeyDevice, device)
			rWithUser := r.WithContext(ctxWithUser)
			next.ServeHTTP(w, rWithUser)
		})
	}
}

func GetAuthenticatedDevice(r *http.Request) (*business.Device, error) {
	token := r.Header.Get("x-schism-token")
	// Get device id from token
	queryStmt := "SELECT device_id FROM accesstokens where id = ?"
	stmt, err := db.DB.Prepare(queryStmt)
	if err != nil {
		return nil, err
	}
	var deviceId string
	err = stmt.QueryRow(token).Scan(&deviceId)
	if err != nil {
		return nil, err
	}
	// Get device
	queryStmt = "SELECT id, name FROM devices where id = ?"
	stmt, err = db.DB.Prepare(queryStmt)
	if err != nil {
		return nil, err
	}
	var device business.Device
	err = stmt.QueryRow(deviceId).Scan(&device.Id, &device.Name)
	if err != nil {
		return nil, err
	}
	return &device, nil
}
