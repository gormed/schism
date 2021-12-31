package permissions

import (
	"net/http"

	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/api/middleware"
	"gitlab.void-ptr.org/go/schism/pkg/business"
)

func HasPermission(w http.ResponseWriter, r *http.Request, deviceId string) bool {
	self := r.Context().Value(middleware.ContextKeyDevice).(*business.Device)
	if *self.Id != deviceId {
		http.Error(w, errors.StatusForbidden, http.StatusForbidden)
		return false
	}
	return true
}
