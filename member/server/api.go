package server

import (
	"github.com/nachtenontij/infra/base/server"
	"github.com/nachtenontij/infra/member"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req member.LoginRequest
	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	e := ByHandle(req.Handle)
	if e == nil {
		http.Error(w, "no such user", 404)
		return
	}

	server.WriteJsonResponse(w, &member.LoginResponse{
		SessionKey: e.NewSession(),
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

}

func EnlistHandler(w http.ResponseWriter, r *http.Request) {
	var req member.EnlistRequest
	if !server.ReadJsonRequest(w, r, &req) {
		return
	}
	server.WriteJsonResponse(w, "hi")
}
