package base

import (
	"gopkg.in/mgo.v2/bson"
)

type HandleOrId struct {
	Handle *string
	Id     *string
}

// Identifier used in the database.
// We encapsulate Mongo's ObjectId such that we can give Id the Ider interface.
type Id bson.ObjectId

// An object with an Id.
type Ider interface {
	Id() Id
}

func (id *Id) Id() Id {
	return *id
}

func (id *Id) String() string {
	return (*bson.ObjectId)(id).Hex()
}

func (id *Id) GetBSON() (interface{}, error) {
	return (*bson.ObjectId)(id), nil
}

func (id *Id) SetBSON(raw bson.Raw) error {
	var val *bson.ObjectId
	if err := raw.Unmarshal(&val); err != nil {
		return err
	}
	if val == nil {
		return bson.SetZero
	}
	*id = *(*Id)(val)
	return nil
}

func (id *Id) MarshalJSON() ([]byte, error) {
	return (*bson.ObjectId)(id).MarshalJSON()
}

func (id *Id) UnmarshalJSON(b []byte) error {
	return (*bson.ObjectId)(id).UnmarshalJSON(b)
}

// Returns a new unique Id.
func NewId() Id {
	return Id(bson.NewObjectId())
}

// Parses and returns an Id in Hex form.  Returns nil if malformed.
func IdHex(hex string) *Id {
	if !bson.IsObjectIdHex(hex) {
		return nil
	}
	ret := Id(bson.ObjectIdHex(hex))
	return &ret
}

// Convert []Iders to []Id
func Ids(iders []Ider) (ret []Id) {
	for _, ider := range iders {
		ret = append(ret, ider.Id())
	}
	return
}
