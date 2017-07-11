// Members, groups (they're all entities!)

package member

import (
	"github.com/nachtenontij/infra/base"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"time"
)

type EntityData struct {
	Id bson.ObjectId `bson:"id"`

	Kind Kind

	Name           string
	GenitivePrefix base.Text

	// Used as loginname and identifiers
	Handles []string

	Group *GroupData
	User  *UserData
}

// Kind of entity
type Kind int8

const (
	Group Kind = iota
	User
)

type UserData struct {
	// Contact
	EMail        string
	Address      *Address
	Phonenumbers []string

	Person Person

	InvitedBy *bson.ObjectId

	// Security
	PasswordHash *string
}

type Person struct {
	GivenNames    []string // e.g. ["Claire", "Marie"]
	TussenVoegsel string   // e.g. van der, van
	LastName      string
	Nickname      string // e.g. Homer
	Gender        string // e.g. "", "vrouw", "man", ...
	Prefix        string // e.g. Mr.
	Suffix        string // e.g. M.Sc

	DateOfBirth time.Time
}

// Address of a user
type Address struct {
	Street  string
	Number  string
	State   string // e.g. Gelderland
	City    string
	Zip     string
	Country string
}

// Group specific data
type GroupData struct {
	Description base.Text
}

// Relation between entities
type RelationData struct {
	Id    bson.ObjectId `bson:"id"`
	From  time.Time
	Until time.Time
	How   *bson.ObjectId
}

type SessionData struct {
	Key          string
	UserId       *bson.ObjectId
	Created      time.Time
	LastActivity time.Time
	IsGenesis    bool
}

// Returns the full name of a person
func (p *Person) Name() string {
	var ret string
	if p.Nickname != "" {
		ret += p.Nickname + " "
	} else if len(p.GivenNames) != 0 {
		ret += p.GivenNames[0] + " "
	}
	if p.TussenVoegsel != "" {
		ret += p.TussenVoegsel + " "
	}
	ret += p.LastName
	return ret
}

// Regular expression used by ValidHandle
var rxHandle = regexp.MustCompile("^[0-9a-z][0-9a-z-]*$")

// Check whether the handle is valid
func ValidHandle(handle string) bool {
	return rxHandle.MatchString(handle)
}
