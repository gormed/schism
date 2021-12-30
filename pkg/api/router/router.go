package router

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	Devices map[string]interface{} `json:"devices"`
}

var routerMap = routesMap{}
var Database *db.Sqlite = nil

func SchismRouter(db *db.Sqlite) *mux.Router {
	if db == nil {
		util.Log.Fatalf("no database given for initialization")
	}

	Database = db
	api.ApiSecret = api.ReadSecret("schism.api.secret")
	r := mux.NewRouter()

	// Setup CORS
	// IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
	r.Use(mux.CORSMethodMiddleware(r))

	// Setup request logging
	logMiddleware := _middleware.NewLogMiddleware(util.Log)
	r.Use(logMiddleware.Func())

	// Create our middlewares
	secretMiddleware := middleware.NewSecretMiddleware(api.ApiSecret)
	authMiddleware := middleware.NewAuthMiddleware()

	if api.Features.Devices.Enabled {
		routerMap.Devices = map[string]interface{}{
			"GET":    "/devices/{id}",
			"POST":   "/devices",
			"PATCH":  "/devices/{id}",
			"DELETE": "/devices/{id}",
		}
		deviceHandler := &handler.DeviceHandler{
			Database: Database,
		}

		// Public device route (POST)
		publicDeviceRouter := r.NewRoute().Subrouter()

		publicDeviceRouter.Use(secretMiddleware.Func())

		publicDeviceRouter.HandleFunc("/devices", deviceHandler.MakePostDevice()).Methods("POST", "OPTIONS")

		// Private device routes (GET, PATCH, DELETE)
		privateDeviceRouter := r.NewRoute().Subrouter()

		privateDeviceRouter.Use(secretMiddleware.Func())
		privateDeviceRouter.Use(authMiddleware.Func())

		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.MakeGetDevice()).Methods("GET", "OPTIONS")
		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.MakePatchDevice()).Methods("PATCH", "OPTIONS")
		privateDeviceRouter.HandleFunc("/devices/{id}", deviceHandler.MakeDeleteDevice()).Methods("DELETE", "OPTIONS")
	}

	// Write out api infos
	r.HandleFunc("/", MakeDefaultHandler(routerMap)).Methods("GET", "OPTIONS")

	// Not found handler
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
