package server

import (
	"net/http"
    "github.com/nachtenontij/infra/member"
    "github.com/nachtenontij/infra/base/server"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req member.LoginRequest
    if !server.ReadJsonRequest(w, r, &req) {
        return
    }
    server.WriteJsonResponse(w, "hi")
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
