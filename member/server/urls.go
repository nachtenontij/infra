package server

import (
	"github.com/gorilla/mux"
)

func RegisterUrls(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()

	s.HandleFunc("/login", LoginHandler).Methods("POST")
	s.HandleFunc("/logout", LogoutHandler).Methods("POST")
	s.HandleFunc("/enlist", EnlistHandler).Methods("POST")
}

