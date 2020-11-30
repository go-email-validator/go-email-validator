package performance

import (
	"reflect"
	"testing"
)

type type1 struct{}
type type2 struct{}
type type3 struct{}
type type4 struct{}
type type5 struct{}
type type6 struct{}
type type7 struct{}
type type8 struct{}
type type9 struct{}
type type10 struct{}
type type11 struct{}
type type12 struct{}
type type13 struct{}
type type14 struct{}
type type15 struct{}
type type16 struct{}
type type17 struct{}
type type18 struct{}
type type19 struct{}
type type20 struct{}

var typeSlice = []reflect.Type{
	reflect.TypeOf((*type1)(nil)),
	reflect.TypeOf((*type2)(nil)),
	reflect.TypeOf((*type3)(nil)),
	reflect.TypeOf((*type4)(nil)),
	reflect.TypeOf((*type5)(nil)),
	reflect.TypeOf((*type6)(nil)),
	reflect.TypeOf((*type7)(nil)),
	reflect.TypeOf((*type8)(nil)),
	reflect.TypeOf((*type9)(nil)),
	reflect.TypeOf((*type10)(nil)),
	reflect.TypeOf((*type11)(nil)),
	reflect.TypeOf((*type12)(nil)),
	reflect.TypeOf((*type13)(nil)),
	reflect.TypeOf((*type14)(nil)),
	reflect.TypeOf((*type15)(nil)),
	reflect.TypeOf((*type16)(nil)),
	reflect.TypeOf((*type17)(nil)),
	reflect.TypeOf((*type18)(nil)),
	reflect.TypeOf((*type19)(nil)),
	reflect.TypeOf((*type20)(nil)),
}

func BenchmarkGetType(b *testing.B) {
	var v reflect.Type
	for i := 1; i < 10000000; i++ {
		v = reflect.TypeOf((*type1)(nil))
	}
	_ = v
}
func BenchmarkGetTypeName(b *testing.B) {
	var v string
	for i := 1; i < 10000000; i++ {
		v = reflect.TypeOf((*type1)(nil)).Name()
	}
	_ = v
}

func BenchmarkTypeMap(b *testing.B) {
	m := getTypeMap()

	for i := 1; i < 100000; i++ {
		for _, value := range typeSlice {
			m[value] = true
		}

		var check interface{}
		for _, value := range typeSlice {
			check = m[value]
		}
		_ = check
	}
}

var typeStringSlice = []string{
	reflect.TypeOf((*type1)(nil)).Elem().Name(),
	reflect.TypeOf((*type2)(nil)).Elem().Name(),
	reflect.TypeOf((*type3)(nil)).Elem().Name(),
	reflect.TypeOf((*type4)(nil)).Elem().Name(),
	reflect.TypeOf((*type5)(nil)).Elem().Name(),
	reflect.TypeOf((*type6)(nil)).Elem().Name(),
	reflect.TypeOf((*type7)(nil)).Elem().Name(),
	reflect.TypeOf((*type8)(nil)).Elem().Name(),
	reflect.TypeOf((*type9)(nil)).Elem().Name(),
	reflect.TypeOf((*type10)(nil)).Elem().Name(),
	reflect.TypeOf((*type11)(nil)).Elem().Name(),
	reflect.TypeOf((*type12)(nil)).Elem().Name(),
	reflect.TypeOf((*type13)(nil)).Elem().Name(),
	reflect.TypeOf((*type14)(nil)).Elem().Name(),
	reflect.TypeOf((*type15)(nil)).Elem().Name(),
	reflect.TypeOf((*type16)(nil)).Elem().Name(),
	reflect.TypeOf((*type17)(nil)).Elem().Name(),
	reflect.TypeOf((*type18)(nil)).Elem().Name(),
	reflect.TypeOf((*type19)(nil)).Elem().Name(),
	reflect.TypeOf((*type20)(nil)).Elem().Name(),
}

func BenchmarkTypeStringMap(b *testing.B) {
	m := getStringMap()

	for i := 1; i < 100000; i++ {
		for _, value := range typeStringSlice {
			m[value] = true
		}

		var check interface{}
		for _, value := range typeStringSlice {
			check = m[value]
		}
		_ = check
	}
}

var strSlice = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
	"16", "17", "18", "19", "20",
}

func BenchmarkStringMap(b *testing.B) {
	m := getStringMap()

	for i := 1; i < 100000; i++ {
		for _, value := range strSlice {
			m[value] = true
		}

		var check interface{}
		for _, value := range strSlice {
			check = m[value]
		}
		_ = check
	}
}

var intSlice = []int{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	16, 17, 18, 19, 20,
}

func BenchmarkIntMap(b *testing.B) {
	m := getIntMap()

	for i := 1; i < 100000; i++ {
		for _, value := range intSlice {
			m[value] = true
		}

		var check interface{}
		for _, value := range intSlice {
			check = m[value]
		}
		_ = check
	}
}

func getTypeMap() map[reflect.Type]interface{} {
	return make(map[reflect.Type]interface{})
}

func getStringMap() map[string]interface{} {
	return make(map[string]interface{})
}

func getIntMap() map[int]interface{} {
	return make(map[int]interface{})
}
