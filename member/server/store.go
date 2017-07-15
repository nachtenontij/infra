package server

import (
	"github.com/nachtenontij/infra/base"
	"github.com/nachtenontij/infra/base/server"
	"github.com/nachtenontij/infra/member"
	"gopkg.in/hlandau/passlib.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

var scol *mgo.Collection // sessions collection
var ecol *mgo.Collection // entities collection
var rcol *mgo.Collection // relations collection
var bcol *mgo.Collection // brands collections

func InitializeCollections(db *mgo.Database) {
	// Get collections
	scol = db.C("sessions")
	ecol = db.C("entities")
	rcol = db.C("relations")
	bcol = db.C("brands")

	// Check/create indices
	if err := ecol.EnsureIndex(mgo.Index{
		Key:    []string{"handles"},
		Unique: true,
		Sparse: true,
	}); err != nil {
		log.Fatalf("EnsureIndex entities.Handles: %s", err)
	}

	if err := bcol.EnsureIndex(mgo.Index{
		Key:    []string{"handle"},
		Unique: true,
	}); err != nil {
		log.Fatalf("EnsureIndex brands.Handle: %s", err)
	}

	if err := scol.EnsureIndex(mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}); err != nil {
		log.Fatalf("EnsureIndex sessions.Key: %s", err)
	}

	if err := scol.EnsureIndex(mgo.Index{
		Key:    []string{"isgenesis"},
		Sparse: true,
	}); err != nil {
		log.Fatalf("EnsureIndex sessions.IsGenesis: %s", err)
	}

	if err := scol.EnsureIndex(mgo.Index{
		Key: []string{"lastactivity"},
	}); err != nil {
		log.Fatalf("EnsureIndex sessions.LastActivity: %s", err)
	}

	if err := scol.EnsureIndex(mgo.Index{
		Key: []string{"userid"},
	}); err != nil {
		log.Fatalf("EnsureIndex sessions.UserId: %s", err)
	}

	// Create genesis session
	log.Printf("Genesis session key: %s", server.Settings.GenesisSessionKey)
	scol.RemoveAll(bson.M{"isgenesis": true})
	if err := scol.Insert(&member.SessionData{
		Key:       server.Settings.GenesisSessionKey,
		IsGenesis: true,
	}); err != nil {
		log.Fatalf("Creating genesis session failed: %s", err)
	}

	// TODO remove old genesis sessions
}

type Entity struct {
	data *member.EntityData
}

type Session struct {
	data *member.SessionData
	user *Entity
}

type Brand struct {
	data *member.BrandData
}

type Relation struct {
	data *member.RelationData
}

// Finds entity by id
func ByIdString(id string) *Entity {
	return ById(bson.ObjectIdHex(id))
}

// Checks whether an object exists and returns parsed objectid
func ExistsByIdString(id string) *bson.ObjectId {
	if !bson.IsObjectIdHex(id) {
		return nil
	}
	ret := bson.ObjectIdHex(id)
	n, err := ecol.Find(bson.M{"id": ret}).Count()
	if err != nil {
		log.Printf("ExistsByIdString: %s", err)
		return nil
	}
	if n == 0 {
		return nil
	}
	return &ret
}

// Finds entity by id
func ById(id bson.ObjectId) *Entity {
	var data member.EntityData
	if ecol.Find(bson.M{"id": id}).One(&data) != nil {
		return nil
	}
	return fromData(&data)
}

// Find Entity by handle
func ByHandle(handle string) *Entity {
	var data member.EntityData
	if ecol.Find(bson.M{"handles": handle}).One(&data) != nil {
		return nil
	}
	return fromData(&data)
}

// Find brand by handle
func BrandByHandle(handle string) *Brand {
	var data member.BrandData
	if ecol.Find(bson.M{"handle": handle}).One(&data) != nil {
		return nil
	}
	return brandFromData(&data)
}

// Find Entity ID by handle
func IdByHandle(handle string) *bson.ObjectId {
	var data member.EntityData
	if ecol.Find(bson.M{"handle": handle}).Select(
		bson.M{"id": 1}).One(&data) != nil {
		return nil
	}
	return &data.Id
}

// Creates an Entity object from an EntityData object
func fromData(data *member.EntityData) *Entity {
	return &Entity{data: data}
}

// Creates a Brand object from an BrandData object
func brandFromData(data *member.BrandData) *Brand {
	return &Brand{data: data}
}

// Finds session by key
func SessionByKey(key string) *Session {
	var data member.SessionData
	if scol.Find(bson.M{"key": key}).One(&data) != nil {
		return nil
	}
	return &Session{data: &data}
}

// Returns session associated to the request
func SessionFromRequest(r *http.Request) *Session {
	session, _ := r.Context().Value("session").(*Session)
	return session
}

// Returns session and user associated to the request.
func SessionUserFromRequest(r *http.Request) (*Session, *Entity) {
	session := SessionFromRequest(r)
	return session, session.User()
}

func (s *Session) User() *Entity {
	if s.user == nil {
		if s.data.UserId == nil {
			return nil
		}
		s.user = ById(*s.data.UserId)
	}
	return s.user
}

// Updates LastActivity on Session
func (s *Session) Touch() {
	s.data.LastActivity = time.Now()
	s.Save()
}

// Saves the session to the database
func (s *Session) Save() {
	if err := scol.Update(bson.M{"key": s.data.Key}, s.data); err != nil {
		log.Printf("Session.Save(): scol.Update(): %s", err)
	}
}

func (s *Session) Logout(all bool) {
	if err := scol.Remove(bson.M{"key": s.data.Key}); err != nil {
		log.Printf("Session.Logout(): %s", err)
	}
	if !all {
		return
	}
	if _, err := scol.RemoveAll(
		bson.M{"userid": s.data.UserId}); err != nil {
		log.Printf("Session.Logout(): All: %s", err)
	}
}

// Checks whether the user associated to the session is allowed to
// view and edit the full membership database
func (s *Session) IsMemberAdmin() bool {
	if s.data.IsGenesis {
		return true
	}
	// TODO
	return false
}

// Saves the entity to the database
func (e *Entity) Save() {
	if err := ecol.Update(bson.M{"id": e.data.Id}, e.data); err != nil {
		log.Printf("Entity.Save(): ecol.Update(): %s", err)
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
	e.AssertUser()
	key := base.GenerateHexSecret(32)
	data := member.SessionData{
		Key:          key,
		UserId:       &e.data.Id,
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

func (e *Entity) SetPassword(password string) (err error) {
	e.AssertUser()

	hash, err := passlib.Hash(password)
	if err != nil {
		return
	}

	e.data.User.PasswordHash = &hash
	go e.Save()

	return nil
}
