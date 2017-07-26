// Members, groups (they're all entities!)

package member

import (
	"github.com/nachtenontij/infra/base"
	"regexp"
	"time"
)

type EntityData struct {
	Id base.Id `bson:"_id"`

	Kind Kind `read:all`

	Name           string    `read:"all"`
	GenitivePrefix base.Text `read:"all"`

	// Used as loginname and identifiers
	Handles []string `read:all`

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
	EMail        string   `read:"user",write:"user"`
	Address      *Address `read:"user",write:"user"`
	Phonenumbers []string `read:"user",write:"user"`

	Person Person

	InvitedBy *base.Id

	// Security
	PasswordHash *string
}

type Person struct {
	// e.g. ["Claire", "Marie"]
	GivenNames []string `read:"all",write:"admin"`
	// e.g. van der, van
	TussenVoegsel string `read:"all",write:"admin"`
	LastName      string `read:"all",write:"admin"`
	// e.g. Homer
	Nickname string `read:"all",write:"admin"`
	// e.g. "", "vrouw", "man", ...
	Gender string `read:"all",write:"user"`
	// e.g. Mr.
	Prefix string `read:"all",write:"user"`
	// e.g. M.Sc
	Suffix      string    `read:"all",write:"user"`
	DateOfBirth time.Time `read:"user",write:"user"`
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

// MinTime represents from-the-start-of-time
var MinTime = time.Date(1903, 12, 28, 0, 0, 0, 0, time.UTC)

// MaxTime represents into-the-foreseeable future
var MaxTime = time.Date(11876, 4, 4, 1, 1, 9, 2, time.UTC)

// Relation between entities
type RelationData struct {
	Id    base.Id `bson:"_id"`
	From  time.Time
	Until time.Time

	// What kind of relation is this?  This the nil for simple membership
	// or the handle of a brand, e.g. "chair".
	How  *string
	Who  base.Id
	With base.Id
}

// A brand represents the type of a relation
type BrandData struct {
	Handle      string
	Name        base.Text
	Description base.Text
}

// Data for a user session.  A user can have multiple sessions.
type SessionData struct {
	Key          string
	UserId       *base.Id
	Created      time.Time
	LastActivity time.Time

	// When start up the ontijd daemon will create a "genisis" admin session
	// (without user attached) to allow scripts to interact (even if
	// there is no user yet).
	IsGenesis bool
}

// Records sensitive actions.  These records are intended to be read
// by humans --- not by scripts.
type AuditRecordData struct {
	Id      base.Id  `bson:"_id"`
	By      *base.Id // Who performed the action.
	Entity  base.Id  // On whom/what was the action performed
	Message string   // What happened
	When    time.Time
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
