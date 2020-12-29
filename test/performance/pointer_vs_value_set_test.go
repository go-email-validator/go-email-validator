package performance

import "testing"

var value = 10
var valueForPtr = value
var ptr = &valueForPtr

func BenchmarkPointerSet(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		*ptr++
		_ = *ptr
	}
}

func BenchmarkValueSet(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		value++
		_ = value
	}
}
