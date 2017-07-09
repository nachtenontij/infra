package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type RegisterUrlsFunc func(*mux.Router)
type MiddlewareFunc func(http.Handler) http.Handler

var registerUrlsFuncs []RegisterUrlsFunc
var middlewareFuncs []MiddlewareFunc

func RegisterRegisterUrls(f RegisterUrlsFunc) {
	registerUrlsFuncs = append(registerUrlsFuncs, f)
}

// TODO priority for middleware?
func RegisterMiddleware(f MiddlewareFunc) {
	middlewareFuncs = append(middlewareFuncs, f)
}

func ListenAndServe() error {
	err := ConnectToDatabase()
	if err != nil {
		return fmt.Errorf("Could not connect to database: %s", err)
	}

	r := mux.NewRouter()
	for _, registerUrls := range registerUrlsFuncs {
		registerUrls(r)
	}

	var handler http.Handler = r
	for _, middleware := range middlewareFuncs {
		handler = middleware(handler)
	}

	http.Handle("/", handler)

	log.Printf("Listening on %s ...", Settings.BindAddress)

	return http.ListenAndServe(Settings.BindAddress, nil)
}

func WriteJsonResponse(w http.ResponseWriter, resp interface{}) {
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ReadJsonRequest(w http.ResponseWriter,
	r *http.Request,
	v interface{}) bool {
	err := json.Unmarshal([]byte(r.FormValue("request")), v)
	if err != nil {
		http.Error(w, fmt.Sprintf(
			"Missing or malformed request parameter: %s", err), 400)
		return false
	}
	return true
}
