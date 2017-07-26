package member

import (
	"github.com/nachtenontij/infra/base"
)

type LoginRequest struct {
	Handle   string
	Password string
}

type LoginResponse struct {
	SessionKey string
}

type LogoutRequest struct {
	// when true log out all sessions of this user
	All bool
}

type EnlistRequest struct {
	Handle      string
	EMail       string
	Address     Address
	Phonenumber string
	InvitedBy   string
	Person      Person
}

type EnlistResponse struct {
	Id string
}

type SelectUserRequest struct {
	Handle string
}

type SelectUserResponse struct {
	Id string
}

type PasswdRequest struct {
	Password string
}

type PasswdResponse struct {
}

type GetEntityRequest struct {
	Which base.HandleOrId
}

type GetEntityResponse struct {
	// partial EntityData
	Entity base.Patch
}

type PatchEntityRequest struct {
	Which base.HandleOrId `json:"inline"`

	// EntityData delta
	Patch base.Patch
}
type PatchEntityResponse struct {
}
