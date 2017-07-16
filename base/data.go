package base

import (
	"reflect"
)

type HandleOrId struct {
	Handle *string
	Id     *string
}

type Patch map[string]interface{}

// Diff returns a patch that when applied to 'from' yields 'to'.
func Diff(to, from interface{}) Patch {
	return PatchFrom(to, func(field reflect.StructField,
		tovalue reflect.Value, others ...reflect.Value) bool {

		return reflect.DeepEqual(tovalue.Interface(),
			others[0].Interface())

	}, []interface{}{from})
}

func PatchFrom(obj interface{},
	filter func(reflect.StructField, reflect.Value, ...reflect.Value) bool,
	others ...interface{}) Patch {

	result := map[string]interface{}{}

	WalkStruct(func(path string, field reflect.StructField,
		values ...reflect.Value) {

		if filter(field, values[0], values[1:]...) {
			return
		}

		result[path] = values[0].Interface()

	}, append([]interface{}{obj}, others...))

	return result
}

func WalkStruct(handler func(string, reflect.StructField, ...reflect.Value),
	structs ...interface{}) {

	walkStruct{handler: handler}.do(structs)
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

func (c walkStruct) do(structs []interface{}) {
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

func (c walkStruct) doOne(task walkStructTask) {
	t := task.values[0].Type()

	switch t.Kind() {
	case reflect.Struct:
		c.doOneStruct(task, t)
	default:
		c.handler(task.path, task.field, task.values...)
	}
}

func (c walkStruct) doOneStruct(task walkStructTask, t reflect.Type) {
	numfield := t.NumField()
	for i := 0; i < numfield; i++ {
		newtask := walkStructTask{}

		newtask.field = t.Field(i)
		newtask.path = task.path + "/" + newtask.field.Name
		newtask.values = make([]reflect.Value, 0, len(task.values))

		for j := 0; j < len(task.values); j++ {
			newtask.values = append(newtask.values,
				task.values[j].Field(i))
		}

		c.todo = append(c.todo, newtask)
	}
}
