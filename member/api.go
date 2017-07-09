package member

import (
    // "github.com/nachtenontij/infra/base"
	// "gopkg.in/mgo.v2/bson"
)

type LoginRequest struct {
    Handle string
    Password string
}

type LoginResponse struct {
    SessionKey string
}

type CommonDataResponse struct {
    PasskeySalt string
}
