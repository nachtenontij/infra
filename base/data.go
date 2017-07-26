package base

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
)

type HandleOrId struct {
	Handle *string
	Id     *string
}

type Patch map[string]interface{}

// Diff returns a patch that when applied to 'from' yields 'to'.
func Diff(to, from interface{}) Patch {
	return PatchFrom(to,
		func(field reflect.StructField, tovalue reflect.Value,
			others ...reflect.Value) bool {

			return !reflect.DeepEqual(tovalue.Interface(),
				others[0].Interface())

		}, []interface{}{from})
}

func PatchFromTags(obj interface{}, filter func(reflect.StructTag) bool) Patch {

	return PatchFrom(obj,
		func(field reflect.StructField, v reflect.Value,
			o ...reflect.Value) bool {

			return filter(field.Tag)
		})
}

func PatchFrom(obj interface{},
	filter func(reflect.StructField, reflect.Value, ...reflect.Value) bool,
	others ...interface{}) Patch {

	result := map[string]interface{}{}

	WalkStruct(func(path string, field reflect.StructField,
		values ...reflect.Value) {

		if !filter(field, values[0], values[1:]...) {
			return
		}

		result[path] = values[0].Interface()

	}, append([]interface{}{obj}, others...)...)

	return result
}

func WalkStruct(handler func(string, reflect.StructField, ...reflect.Value),
	structs ...interface{}) {

	(&walkStruct{handler: handler}).do(structs)
}

type walkStruct struct {
	handler func(string, reflect.StructField, ...reflect.Value)
	todo    []walkStructTask
}

type walkStructTask struct {
	path   string
	field  reflect.StructField
	values []reflect.Value
}

func (c *walkStruct) do(structs []interface{}) {
	if len(structs) == 0 {
		panic("zero structs given")
	}

	values := make([]reflect.Value, 0, len(structs))
	for i := 0; i < len(structs); i++ {
		values = append(values, reflect.ValueOf(structs[i]))
	}

	task := walkStructTask{path: "", values: values}

	c.todo = []walkStructTask{task}

	for len(c.todo) > 0 {
		task, c.todo = c.todo[0], c.todo[1:]
		c.doOne(task)
	}
}

func (c *walkStruct) doOne(task walkStructTask) {
	t := task.values[0].Type()

	switch t.Kind() {
	case reflect.Ptr:
		c.doOnePtr(task, t)
	case reflect.Struct:
		c.doOneStruct(task, t)
	default:
		c.handler(task.path, task.field, task.values...)
	}
}

func (c *walkStruct) doOneStruct(task walkStructTask, t reflect.Type) {
	numfield := t.NumField()
	for i := 0; i < numfield; i++ {
		f := t.Field(i)

		if f.PkgPath != "" {
			// unexported field
			continue
		}

		newtask := walkStructTask{}

		newtask.field = f
		newtask.path = task.path + "." + newtask.field.Name
		newtask.values = make([]reflect.Value, 0, len(task.values))

		for j := 0; j < len(task.values); j++ {
			v := task.values[j]

			if v.Type() != t {
				log.Panicf("WalkStruct: "+
					"structs not of the same type %+v",
					task)
			}

			newtask.values = append(newtask.values, v.Field(i))
		}

		c.todo = append(c.todo, newtask)
	}
}

func (c *walkStruct) doOnePtr(task walkStructTask, t reflect.Type) {
	newtask := walkStructTask{}

	newtask.field = task.field
	newtask.path = task.path + "*"
	newtask.values = make([]reflect.Value, 0, len(task.values))

	for j := 0; j < len(task.values); j++ {
		v := task.values[j]
		if v.IsNil() {
			c.handler(task.path, task.field, task.values...)
			return
		}
		newtask.values = append(newtask.values, v.Elem())
	}

	c.todo = append(c.todo, newtask)
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
