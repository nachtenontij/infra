package member

import (
	"gopkg.in/mgo.v2/bson"
)

type LoginRequest struct {
	Handle   string
	Password string
}

type LoginResponse struct {
	SessionKey string
}

type CommonDataResponse struct {
	PasskeySalt string
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
	Id bson.ObjectId
}
