package server

import (
	"fmt"
	"github.com/gorilla/mux"
	memberServer "github.com/nachtenontij/infra/member/server"
	"net/http"
)

func ListenAndServe() error {
	err := ConnectToDatabase()
	if err != nil {
		return fmt.Errorf("Could not connect to database: %s", err)
	}

	r := mux.NewRouter()
	memberServer.RegisterUrls(r)
	http.Handle("/", r)
	return http.ListenAndServe(Settings.BindAddress, nil)
}
