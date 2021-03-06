package server

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/nachtenontij/infra/base"
	"github.com/nachtenontij/infra/base/server"
	"github.com/nachtenontij/infra/member"
	"net/http"
	"reflect"
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
	var req member.LogoutRequest
	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	session := SessionFromRequest(r)
	if session != nil {
		session.Logout(req.All)
	}
	server.WriteJsonResponse(w, session != nil)
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

	resp.Id = user.Id()
	session.data.UserId = &resp.Id

	go session.Save()

	server.WriteJsonResponse(w, &resp)
}

func PasswdHandler(w http.ResponseWriter, r *http.Request) {
	var req member.PasswdRequest
	var resp member.PasswdResponse
	session, user := SessionUserFromRequest(r)

	if session == nil {
		http.Error(w, "access denied", 403)
		return
	}

	if user == nil {
		http.Error(w, "no user in session", 400)
		return
	}

	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	if err := user.SetPassword(req.Password); err != nil {
		http.Error(w, "failed to hash password", 400)
		return
	}

	user.AuditLog(user, "Changed password")

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
	resp.Id = base.NewId()

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

	if err := AddEntity(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to insert: %s", err), 400)
		return
	}

	server.WriteJsonResponse(w, &resp)
}

func GetEntityHandler(w http.ResponseWriter, r *http.Request) {
	var req member.GetEntityRequest
	var resp member.GetEntityResponse
	session, user := SessionUserFromRequest(r)

	if session == nil {
		http.Error(w, "access denied", 403)
		return
	}

	if !server.ReadJsonRequest(w, r, &req) {
		return
	}

	e := ByHoi(req.Which)
	if e == nil {
		http.Error(w, "no such entity", 404)
		return
	}

	clearance := map[string]bool{
		"all":   true,
		"user":  false,
		"admin": false,
	}

	if user != nil && user.data.Id == e.data.Id {
		clearance["user"] = true
	}

	if session.IsMemberAdmin() {
		clearance["admin"] = true
		clearance["user"] = true
	}

	resp.Entity = base.PatchFromTags(e.data,
		func(tag reflect.StructTag) bool {
			var v string

			v, ok := tag.Lookup("read")
			if !ok {
				return false
			}

			access, ok := clearance[v]

			if !ok {
				panic("unknown confidentiality level")
			}

			return access
		})

	server.WriteJsonResponse(w, &resp)
}
