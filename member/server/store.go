package server

import (
	"github.com/nachtenontij/infra/member"
    "gopkg.in/hlandau/passlib.v1"
	"gopkg.in/mgo.v2"
)

var scol *mgo.Collection
var ecol *mgo.Collection
var rcol *mgo.Collection

func InitializeCollections(db *mgo.Database) {
    scol = db.C("sessions")
    ecol = db.C("entities")
    rcol = db.C("relations")
}

type Entity struct {
	data *member.EntityData
}

type Session struct {
	data *member.SessionData
}

func (e *Entity) CheckPassword(password string) bool {
    if e.data.User == nil || e.data.User.PasswordHash == nil {
        return false
    }
    _, err := passlib.Verify(password, *e.data.User.PasswordHash)
    if err != nil {
        return false
    }
    return true
}
