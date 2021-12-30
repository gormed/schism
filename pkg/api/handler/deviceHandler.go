package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.void-ptr.org/go/schism/pkg/api/middleware"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DeviceRequest struct {
}

// MakeGetDevice ...
func MakeGetDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		device := r.Context().Value(middleware.ContextKeyDevice).(*business.Device)
		err := json.NewEncoder(w).Encode(device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}

// MakePostDevice ...
func MakePostDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var request DeviceRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			util.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// MakePatchDevice ...
func MakePatchDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var request DeviceRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			util.Log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// MakeDeleteDevice ...
func MakeDeleteDevice() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(http.StatusOK)
	}
}
