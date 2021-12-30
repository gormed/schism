package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/api/middleware"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
)

type DeviceRequest struct {
	business.Device
}

type DeviceHandler struct {
	Database *db.Sqlite
}

func (dh *DeviceHandler) hasPermission(w http.ResponseWriter, r *http.Request, deviceId string) bool {
	self := r.Context().Value(middleware.ContextKeyDevice).(*business.Device)
	if *self.Id != deviceId {
		http.Error(w, errors.StatusForbidden, http.StatusForbidden)
		return false
	}
	return true
}

// MakeGetDevice ...
func (dh *DeviceHandler) MakeGetDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		deviceId := mux.Vars(r)["id"]
		if !dh.hasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		device := &business.Device{Identifyable: db.Identifyable{Id: &deviceId, Database: dh.Database}}
		device, status, err := device.Read()
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// MakePostDevice ...
func (dh *DeviceHandler) MakePostDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var deviceCreate business.DeviceCreate
		err := json.NewDecoder(r.Body).Decode(&deviceCreate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device := &business.Device{Identifyable: db.Identifyable{Database: dh.Database}}
		device, status, err := device.Create(&deviceCreate)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// MakePatchDevice ...
func (dh *DeviceHandler) MakePatchDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		deviceId := mux.Vars(r)["id"]
		if !dh.hasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		var deviceUpdate business.DeviceUpdate
		err := json.NewDecoder(r.Body).Decode(&deviceUpdate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device := &business.Device{Identifyable: db.Identifyable{Id: &deviceId, Database: dh.Database}}
		device, status, err := device.Update(&deviceUpdate)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// MakeDeleteDevice ...
func (dh *DeviceHandler) MakeDeleteDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		deviceId := mux.Vars(r)["id"]
		if !dh.hasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		device := &business.Device{Identifyable: db.Identifyable{Id: &deviceId, Database: dh.Database}}
		device, status, err := device.Delete()
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
