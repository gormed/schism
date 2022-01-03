package permissions

import (
	"net/http"

	"gitlab.void-ptr.org/go/schism/pkg/api"
	"gitlab.void-ptr.org/go/schism/pkg/business"
)

func HasPermission(w http.ResponseWriter, r *http.Request, deviceId string) bool {
	self := r.Context().Value(api.ContextKeyDevice).(*business.Device)
	if *self.Id != deviceId {
		return false
	}
	return true
}
