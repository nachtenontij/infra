package base_test

import (
	"fmt"
	"github.com/nachtenontij/infra/base"
	"reflect"
)

type TestStruct struct {
	S string      `ok:"false"`
	T *TestStruct `ok:"true"`
}

func ExampleWalkStruct() {
	base.WalkStruct(
		func(p string, f reflect.StructField, values ...reflect.Value) {
			fmt.Println(p)
		}, TestStruct{T: &TestStruct{}})
	// Output:
	// .S
	// .T*.S
	// .T*.T
}

func ExamplePatchFromTags() {
	fmt.Println(base.PatchFromTags(TestStruct{T: &TestStruct{}},
		func(tag reflect.StructTag) bool {
			return tag.Get("ok") == "true"
		}))
	// Output:
	// map[.T*.T:<nil>]
}
