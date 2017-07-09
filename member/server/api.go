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

	if !e.CheckPassword(req.Password) {
		http.Error(w, "bad password", 403)
		return
	}

	server.WriteJsonResponse(w, &member.LoginResponse{
		SessionKey: e.NewSession(),
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	SessionFromRequest(r).Logout()
	server.WriteJsonResponse(w, true)
}

func EnlistHandler(w http.ResponseWriter, r *http.Request) {
	var req member.EnlistRequest
	session := SessionFromRequest(r)
	if !session.IsMemberAdmin() {
		http.Error(w, "access denied", 403)
		return
	}
	if !server.ReadJsonRequest(w, r, &req) {
		return
	}
}
