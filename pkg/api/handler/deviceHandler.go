package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/api/headers"
	"gitlab.void-ptr.org/go/schism/pkg/api/permissions"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DeviceRequest struct {
	business.Device
}

type DeviceHandler struct {
	Database *db.Sqlite `json:"-"`
}

// ReadDevice ...
func (dh *DeviceHandler) ReadDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["id"]

		if !permissions.HasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}
		device := business.NewDevice(&deviceId, dh.Database)
		device, status, err := device.Read()
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// CreateDevice ...
func (dh *DeviceHandler) CreateDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var deviceCreate business.DeviceCreate
		err := json.NewDecoder(r.Body).Decode(&deviceCreate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device := business.NewDevice(nil, dh.Database)
		device, status, err := device.Create(&deviceCreate)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// UpdateDevice ...
func (dh *DeviceHandler) UpdateDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["id"]

		if !permissions.HasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		var deviceUpdate business.DeviceUpdate
		err := json.NewDecoder(r.Body).Decode(&deviceUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		device := business.NewDevice(&deviceId, dh.Database)
		device, status, err := device.Update(&deviceUpdate)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// DeleteDevice ...
func (dh *DeviceHandler) DeleteDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["id"]

		if !permissions.HasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}
		device := business.NewDevice(&deviceId, dh.Database)
		device, status, err := device.Delete()
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// LoginDevice ...
func (dh *DeviceHandler) LoginDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["id"]

		// Get secret header
		secret := r.Header.Get(headers.HeaderSchismSecret)

		// Check for the valid api secret
		if api.ApiSecret != secret {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		accesstoken := business.NewAccesstoken(nil, dh.Database)
		accesstoken, status, err := accesstoken.Create(&business.AccesstokenCreate{DeviceId: deviceId})
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(accesstoken)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// LogoutDevice ...
func (dh *DeviceHandler) LogoutDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["id"]

		if !permissions.HasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		accesstoken := r.Context().Value(api.ContextKeyToken).(*business.Accesstoken)
		accesstoken, status, err := accesstoken.Delete()
		if err != nil {
			util.Log.Panic(err.Error())
		}

		token := ""
		accesstoken.Token = &token
		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(accesstoken)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}
