package server

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/nachtenontij/infra/member"
	"gopkg.in/hlandau/passlib.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
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

// Finds entity by id
func ByIdString(id string) *Entity {
	return ById(bson.ObjectIdHex(id))
}

// Finds entity by id
func ById(id bson.ObjectId) *Entity {
	var data member.EntityData
	if ecol.Find(bson.M{"_id": id}).One(&data) != nil {
		return nil
	}
	return fromData(&data)
}

// Find Entity by handle
func ByHandle(handle string) *Entity {
	var data member.EntityData
	if ecol.Find(bson.M{"Handle": handle}).One(&data) != nil {
		return nil
	}
	return fromData(&data)
}

// Creates an Entity object from an EntityData object
func fromData(data *member.EntityData) *Entity {
	return &Entity{data: data}
}

// Finds session by key
func SessionByKey(key string) *Session {
	var data member.SessionData
	if scol.Find(bson.M{"Key": key}).One(&data) != nil {
		return nil
	}
	return &Session{data: &data}
}

// Updates LastActivity on Session
func (s Session) Touch() {
	s.data.LastActivity = time.Now()
	s.Save()
}

// Saves the session to the database
func (s Session) Save() {
	if err := scol.Update(bson.M{"Key": s.data.Key}, s.data); err != nil {
		log.Printf("Session.Save(): scol.Update(): %s", err)
	}
}

// Panics if the entity is not a user.
func (e *Entity) AssertUser() {
	if e.data.Kind != member.User {
		panic("AssertUser")
	}
}

// Creates a new session and returns a sessionkey
func (e *Entity) NewSession() string {
	var rawKey [32]byte
	e.AssertUser()
	// generate a session key
	_, err := rand.Read(rawKey[:])
	if err != nil {
		panic("no random")
	}
	key := hex.Dump(rawKey[:])
	data := member.SessionData{
		Key:          key,
		UserId:       e.data.Id,
		Created:      time.Now(),
		LastActivity: time.Now(),
	}
	if err := scol.Insert(data); err != nil {
		log.Printf("NewSession(): scol.Insert(): %s", err)
	}
	return key
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
