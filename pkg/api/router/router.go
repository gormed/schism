package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_middleware "gitlab.void-ptr.org/go/reflection/pkg/api/middleware"
	"gitlab.void-ptr.org/go/schism/pkg/api"
	"gitlab.void-ptr.org/go/schism/pkg/api/handler"
	"gitlab.void-ptr.org/go/schism/pkg/api/meta"
	"gitlab.void-ptr.org/go/schism/pkg/api/middleware"
	"gitlab.void-ptr.org/go/schism/pkg/db"
	"gitlab.void-ptr.org/go/schism/pkg/util"
)

type routesMap struct {
	Devices map[string][]string `json:"devices"`
	Data    map[string][]string `json:"data"`
}

var routerMap = routesMap{}

func SchismRouter(sqlite *db.Sqlite, influxdb *db.Influx) *mux.Router {
	if sqlite == nil {
		util.Log.Fatalf("no database given for initialization")
	}

	api.ApiSecret = api.ReadSecret("schism.api.secret")
	r := mux.NewRouter()

	// Setup CORS
	// IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
	r.Use(mux.CORSMethodMiddleware(r))

	// Setup request logging
	logMiddleware := _middleware.NewLogMiddleware(util.Log)
	r.Use(logMiddleware.Func())

	// Setup timeout handling
	timeOutMiddleware := middleware.NewTimeOutMiddleware(util.Log, 5*time.Second)
	r.Use(timeOutMiddleware.Func())

	// Create our middlewares
	secretMiddleware := middleware.NewSecretMiddleware(api.ApiSecret)
	authMiddleware := middleware.NewAuthMiddleware(sqlite)

	if api.Features.Devices.Enabled {
		routerMap.Devices = map[string][]string{
			"/devices":            {"POST"},
			"/devices/{id}":       {"GET", "PATCH", "DELETE"},
			"/devices/{id}/login": {"POST"},
		}
		deviceHandler := &handler.DeviceHandler{Database: sqlite}

		// Public device route (POST)
		publicDeviceRouter := r.NewRoute().Subrouter()

		publicDeviceRouter.Use(secretMiddleware.Func())

		publicDeviceRouter.HandleFunc("/devices", deviceHandler.CreateDevice()).Methods("POST", "OPTIONS")
		publicDeviceRouter.HandleFunc("/devices/{id}/login", deviceHandler.LoginDevice()).Methods("POST", "OPTIONS")

		// Private device routes (GET, PATCH, DELETE)
		privateDeviceRouter := r.NewRoute().Subrouter()

		privateDeviceRouter.Use(secretMiddleware.Func())
		privateDeviceRouter.Use(authMiddleware.Func())

		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.ReadDevice()).Methods("GET", "OPTIONS")
		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.UpdateDevice()).Methods("PATCH", "OPTIONS")
		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.DeleteDevice()).Methods("DELETE", "OPTIONS")
		publicDeviceRouter.HandleFunc("/devices/{id}/logout", deviceHandler.LogoutDevice()).Methods("POST", "OPTIONS")
	}

	if api.Features.Data.Enabled {
		routerMap.Data = map[string][]string{
			"/data":          {"POST"},
			"/data/{source}": {"GET"},
		}
		dataHandler := &handler.DataHandler{Database: influxdb}
		privateDataRouter := r.NewRoute().Subrouter()

		privateDataRouter.Use(secretMiddleware.Func())
		privateDataRouter.Use(authMiddleware.Func())

		privateDataRouter.HandleFunc("/data", dataHandler.CreateData()).Methods("POST", "OPTIONS")
		privateDataRouter.HandleFunc("/data/{deviceId}/{source}", dataHandler.ReadData()).Methods("GET", "OPTIONS")
	}

	// Write out api infos
	r.HandleFunc("/", MakeDefaultHandler(routerMap)).Methods("GET", "OPTIONS")

	// Not found handler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.Log.Info(r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	// Method not allowed handler
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return r
}

// MakeDefaultHandler handles the / route to display api informations
func MakeDefaultHandler(routes routesMap) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		query := r.URL.Query()
		responseType := query.Get("type")

		data := map[string]interface{}{
			"meta":     meta.MetaInfo,
			"features": api.Features,
			"routes":   routes,
		}

		switch responseType {
		case "html":
			// w.WriteHeader(http.StatusOK)
			// err := templates["index"].Execute(w, data)
			// if err != nil {
			// 	util.Log.Error(err.Error())
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }
		case "json":
		default:
			json, err := json.Marshal(data)
			if err != nil {
				util.Log.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%s", string(json))
		}
	}
}
