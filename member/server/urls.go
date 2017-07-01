package server

import (
	"github.com/gorilla/mux"
)

func RegisterUrls(r *mux.Router) {
	s := r.PathPrefix("/api").Methods("POST").Subrouter()
	s.HandleFunc("/login", LoginHandler)
}
