// Members, groups (they're all entities!)

package member

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type EntityData struct {
	Id bson.ObjectId `bson:"id"`

	Kind Kind

	Name           string
	GenitivePrefix Text

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
	EMail   string
	Address *Address

	// Person
	GivenNames    []string // e.g. ["Claire", "Marie"]
	TussenVoegsel string   // e.g. van der, van
	LastName      string
	Nickname      string // e.g. Homer
	Gender        string // e.g. "", "vrouw", "man", ...
	Prefix        string // e.g. Mr.
	Suffix        string // e.g. M.Sc

	DateOfBirth time.Time
}

type Address struct {
	Street  string
	Number  string
	State   string // e.g. Gelderland
	City    string
	Zip     string
	Country string
}

type GroupData struct {
	Description Text
}

// TODO move
// String in different languages, e.g.
//   {"nl": "groepen",
//    "en": "groups"}
type Text map[string]string
