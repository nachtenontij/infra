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
	Id base.Id
}

type SelectUserRequest struct {
	Handle string
}

type SelectUserResponse struct {
	Id base.Id // TODO change to string
}

type PasswdRequest struct {
	Password string
}

type PasswdResponse struct {
}

// TODO move to base
type HandleOrId struct {
	Handle *string
	Id     *string
}

type GetEntityRequest struct {
	Which HandleOrId `json:"inline"`
}

type GetEntityRespone struct {
	// partial EntityData
	Entity map[string]interface{}
}

type PatchEntityRequest struct {
	Which HandleOrId `json:"inline"`

	// EntityData delta
	Patch map[string]interface{}
}
type PatchEntityResponse struct {
}
