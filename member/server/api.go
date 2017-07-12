package server

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/nachtenontij/infra/base/server"
	"github.com/nachtenontij/infra/member"
	"gopkg.in/mgo.v2/bson"
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

func SelectUserHandler(w http.ResponseWriter, r *http.Request) {
	var req member.SelectUserRequest
	var resp member.SelectUserResponse
	session := SessionFromRequest(r)

	if session == nil || !session.data.IsGenesis {
		http.Error(w, "access denied", 403)
		return
	}

	if session.data.UserId != nil {
		http.Error(w, "user already set", 400)
		return
	}

	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	user := ByHandle(req.Handle)
	if user == nil {
		http.Error(w, "no such user", 400)
		return
	}

	resp.Id = user.data.Id
	session.data.UserId = &user.data.Id
	go session.Save()

	server.WriteJsonResponse(w, &resp)
}

func EnlistHandler(w http.ResponseWriter, r *http.Request) {
	var req member.EnlistRequest
	var resp member.EnlistResponse
	session := SessionFromRequest(r)

	if session == nil || !session.IsMemberAdmin() {
		http.Error(w, "access denied", 403)
		return
	}
	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	if !govalidator.IsEmail(req.EMail) {
		http.Error(w, "malformed email", 400)
		return
	}

	if !member.ValidHandle(req.Handle) {
		http.Error(w, "malformed handle", 400)
		return
	}

	// TODO set genitive prefix
	resp.Id = bson.NewObjectId()
	data := member.EntityData{
		Id:      resp.Id,
		Kind:    member.User,
		Name:    req.Person.Name(),
		Handles: []string{req.Handle},
		User: &member.UserData{
			EMail:   req.EMail,
			Address: &req.Address,
			Person:  req.Person,
		},
	}

	if req.InvitedBy != "" {
		invitedBy := ExistsByIdString(req.InvitedBy)
		if invitedBy == nil {
			http.Error(w, "no such or malformed InvitedBy", 400)
			return
		}
		data.User.InvitedBy = invitedBy
	}

	if req.Phonenumber != "" {
		data.User.Phonenumbers = []string{req.Phonenumber}
	}

	if err := ecol.Insert(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to insert: %s", err), 400)
		return
	}

	server.WriteJsonResponse(w, &resp)
}
