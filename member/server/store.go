package server

import (
	"fmt"
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
var acol *mgo.Collection // auditRecords collections

func InitializeCollections(db *mgo.Database) {
	// Get collections
	scol = db.C("sessions")
	ecol = db.C("entities")
	rcol = db.C("relations")
	bcol = db.C("brands")
	acol = db.C("auditRecords")

	// Check/create indices
	if err := ecol.EnsureIndex(mgo.Index{
		Key:    []string{"handles"},
		Unique: true,
		Sparse: true,
	}); err != nil {
		log.Fatalf("EnsureIndex entities.Handles: %s", err)
	}

	if err := rcol.EnsureIndex(mgo.Index{
		Key:    []string{"how"},
		Sparse: true,
	}); err != nil {
		log.Fatalf("EnsureIndex relations.How: %s", err)
	}

	if err := rcol.EnsureIndex(mgo.Index{
		Key: []string{"with"},
	}); err != nil {
		log.Fatalf("EnsureIndex relations.With: %s", err)
	}

	if err := rcol.EnsureIndex(mgo.Index{
		Key: []string{"who"},
	}); err != nil {
		log.Fatalf("EnsureIndex relations.Who: %s", err)
	}

	if err := rcol.EnsureIndex(mgo.Index{
		Key: []string{"until", "-from"},
	}); err != nil {
		log.Fatalf("EnsureIndex relations.Until/From: %s", err)
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

	if err := acol.EnsureIndex(mgo.Index{
		Key: []string{"by"},
	}); err != nil {
		log.Fatalf("EnsureIndex auditRecord.By: %s", err)
	}

	if err := acol.EnsureIndex(mgo.Index{
		Key: []string{"entity"},
	}); err != nil {
		log.Fatalf("EnsureIndex auditRecords.Entity: %s", err)
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

type AuditRecord struct {
	data *member.AuditRecordData
}

type Relation struct {
	data *member.RelationData
}

// Convenience struct for AddRelation
type NewRelationData struct {
	From  *time.Time
	Until *time.Time
	How   *string
	Who   base.Ider
	With  base.Ider
}

// A query for relations.  Used in the QueryRelations function.
// Each of the fields represents a possible restriction on the relations
// that should be returned.
type RelationQuery struct {
	// Match relation if the Who is one of these.  Ignored if the list is
	// empty.  List should not contain nil.
	Who  []base.Ider
	With []base.Ider // Match relation if With is one of these.  See Who.
	How  []*string   // Match relation if How is one of these. See Who.

	// Match relations whose running time has non-empty intersection
	// with the interval From and Until forms.  A nil value is taken
	// to be ancient history for From and distant future for Until.
	From  *time.Time
	Until *time.Time // See From.
}

// Finds entity by id
func ByIdString(id string) *Entity {
	parsedId := base.IdHex(id)
	if parsedId == nil {
		return nil
	}
	return ById(parsedId)
}

// Checks whether an object exists and returns parsed objectid
func ExistsByIdString(id string) *base.Id {
	ret := base.IdHex(id)
	if ret == nil {
		return nil
	}
	n, err := ecol.Find(bson.M{"_id": ret}).Count()
	if err != nil {
		log.Printf("ExistsByIdString: %s", err)
		return nil
	}
	if n == 0 {
		return nil
	}
	return ret
}

// Finds entity by id
func ById(id base.Ider) *Entity {
	var data member.EntityData
	if ecol.Find(bson.M{"_id": id.Id()}).One(&data) != nil {
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
func IdByHandle(handle string) *base.Id {
	var data member.EntityData
	if ecol.Find(bson.M{"handle": handle}).Select(
		bson.M{"_id": 1}).One(&data) != nil {
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
		s.user = ById(s.data.UserId)
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
	if err := ecol.Update(bson.M{"_id": e.Id()}, e.data); err != nil {
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
	userId := e.Id()
	data := member.SessionData{
		Key:          key,
		UserId:       &userId,
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

// Create and stores a new audit record
func (e *Entity) AuditLog(by *Entity, what string, args ...interface{}) {
	data := member.AuditRecordData{
		Entity:  e.Id(),
		Message: fmt.Sprintf(what, args...),
		When:    time.Now(),
	}
	if by != nil {
		byId := by.Id()
		data.By = &byId
	}
	if err := acol.Insert(&data); err != nil {
		log.Printf("AuditLog(): acol.Insert(): %s", err)
	}
}

func (e *Entity) Id() base.Id {
	return e.data.Id
}

func (r *Relation) Id() base.Id {
	return r.data.Id
}

// Add a relation.
func AddRelation(rel NewRelationData) error {
	var data member.RelationData
	if rel.Who == nil {
		return fmt.Errorf("rel.Who can't be nil")
	}
	if rel.With == nil {
		return fmt.Errorf("rel.With can't be nil")
	}
	data.Who = rel.Who.Id()
	data.With = rel.With.Id()
	data.How = rel.How
	data.Id = base.NewId()
	if rel.From == nil {
		data.From = member.MinTime
	} else {
		data.From = *rel.From
	}
	if rel.Until == nil {
		data.Until = member.MaxTime
	} else {
		data.Until = *rel.Until
	}
	err := rcol.Insert(data)
	if err != nil {
		log.Printf("AddRelation: %s", err)
	}
	return err
}

// Add an entity
func AddEntity(data member.EntityData) error {
	err := ecol.Insert(data)
	if err != nil {
		log.Printf("AddEntity: %s", err)
	}
	return err
}

// Find relations that match any of the given queries.  See RelationQuery
func QueryRelations(queries []RelationQuery) []member.RelationData {
	if len(queries) == 0 {
		log.Printf("Warning: QueryRelations called with empty queries variable")
		return nil
	}

	// We will construct a mongo query of the form
	//      {"$or": bits}
	var bits []map[string]interface{}

	for _, query := range queries {
		bit := make(map[string]interface{})

		// Put Who into mongo query
		if len(query.Who) == 1 {
			bit["who"] = query.Who[0].Id()
		} else if len(query.Who) > 1 {
			bit["who"] = bson.M{"$in": base.Ids(query.Who)}
		}

		// Put With into mongo query
		if len(query.With) == 1 {
			bit["with"] = query.With[0].Id()
		} else if len(query.With) > 1 {
			bit["with"] = bson.M{"$in": base.Ids(query.With)}
		}

		// Put How into mongo query
		if len(query.How) == 1 {
			bit["how"] = query.How[0]
		} else if len(query.How) > 1 {
			bit["how"] = bson.M{"$in": query.How}
		}

		// Handle nil From and Until
		from := member.MinTime
		until := member.MaxTime
		if query.From != nil {
			from = *query.From
		}
		if query.Until != nil {
			bit["until"] = *query.Until
		}

		// For the situation MinTime < From < Until < MaxTime, we have to
		// put in most effort.  In the other cases, we can simplify the
		// bits required for this Query.
		if from == member.MinTime && until == member.MaxTime {
			bits = append(bits, bit)
		} else if from == until ||
			from == member.MinTime ||
			until == member.MaxTime {
			bit["from"] = bson.M{"$lte": until}
			bit["until"] = bson.M{"$gte": from}
			bits = append(bits, bit)
		} else {
			// the tedious case: we need three bits for this query, for
			// three cases.
			var bit_a, bit_b, bit_c map[string]interface{}

			for k, v := range bit {
				bit_a[k] = v
				bit_b[k] = v
				bit_c[k] = v
			}

			// Case a:   query:                |-------|
			//           relation(s):      |-------|
			//                                   |-|
			bit_a["from"] = bson.M{"$gte": member.MinTime} // to hit index
			bit_a["until"] = bson.M{"$lte": until, "$gte": from}

			// Case b:   query:                |-------|
			//           relation(s):               |-------|
			//                                      |-|
			//  (this overlaps with case a, but that doesn't hurt.)
			bit_b["from"] = bson.M{"$lte": until, "$gte": from}
			bit_b["until"] = bson.M{"$gte": member.MinTime} // to hit index

			// Case c:   query:                |-------|
			//           relation:          |--------------|
			bit_c["from"] = bson.M{"$lte": from}
			bit_c["until"] = bson.M{"$gte": until}

			bits = append(bits, bit_a, bit_b, bit_c)
		}
	}

	var finalQuery map[string]interface{}

	// If bits is a one-element array, the $or query does not return anything,
	// even though it should.  Bug in MongoDB?
	if len(bits) == 1 {
		finalQuery = bits[0]
	} else {
		finalQuery = bson.M{"$or": bits}
	}

	// Perform the query!
	cursor := rcol.Find(finalQuery)

	var result []member.RelationData
	if err := cursor.All(&result); err != nil {
		log.Printf("QueryRelations rcol.Find(): %s", err)
		return nil
	}

	return result
}
