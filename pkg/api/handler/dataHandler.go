package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/api/permissions"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type DataRequest struct {
	business.Data
}

type DataHandler struct {
	Database *db.Influx `json:"-"`
}

// CreateData ...
func (dh *DataHandler) CreateData() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var dataCreate business.DataCreate
		err := json.NewDecoder(r.Body).Decode(&dataCreate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := business.NewData(dh.Database)
		data, status, err := data.Create(&dataCreate)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}

// ReadData ...
func (dh *DataHandler) ReadData() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Get parameters
		deviceId := mux.Vars(r)["deviceId"]
		source := mux.Vars(r)["source"]

		if !permissions.HasPermission(w, r, deviceId) {
			http.Error(w, errors.StatusForbidden, http.StatusForbidden)
			return
		}

		var read = &business.DataRead{}
		start := r.URL.Query().Get("start")
		if len(start) < 1 {
			startTime, err := time.Parse(db.DateLayout, start)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			read.Start = startTime
		} else {
			read.Start = time.Now()
		}
		end := r.URL.Query().Get("end")
		if len(end) < 1 {
			endTime, err := time.Parse(db.DateLayout, end)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			read.End = endTime
		} else {
			read.End = time.Now()
		}

		data := business.NewData(dh.Database)
		data.DeviceId = deviceId
		data.Source = source
		data, status, err := data.Read(read)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			util.Log.Panic(err.Error())
		}
	}
}
