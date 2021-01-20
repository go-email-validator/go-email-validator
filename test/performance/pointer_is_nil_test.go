package performance

import (
	"github.com/modern-go/reflect2"
	"reflect"
	"testing"
)

func getError() error {
	return &MyError{}
}

type MyError struct{}

func (e *MyError) Error() string {
	return "MyError"
}
func IsNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func BenchmarkPointerIsNilReflect(b *testing.B) {
	err := getError()

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		res := IsNil(err)
		_ = res
	}
}

func BenchmarkPointerIsNilReflect2(b *testing.B) {
	err := getError()

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		res := reflect2.IsNil(err)
		_ = res
	}
}

func BenchmarkPointerIsNilDefault(b *testing.B) {
	err := getError()

	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		res := err == (*MyError)(nil)
		_ = res
	}
}
