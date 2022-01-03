package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	_business "gitlab.void-ptr.org/go/reflection/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/api/errors"
	"gitlab.void-ptr.org/go/schism/pkg/api/permissions"
	"gitlab.void-ptr.org/go/schism/pkg/business"
	"gitlab.void-ptr.org/go/schism/pkg/db"
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

		var dataCreate _business.DataCreate
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
			panic(err)
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
		read.Start = r.URL.Query().Get("start")
		read.Stop = r.URL.Query().Get("stop")

		data := business.NewData(dh.Database)
		data.DeviceId = deviceId
		data.Source = source
		result, status, err := data.Read(read)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.WriteHeader(status)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			panic(err)
		}
	}
}
