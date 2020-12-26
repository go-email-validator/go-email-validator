package performance

import "testing"

func BenchmarkRemoveLastElement(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		slice := []string{"A", "B", "C", "D", "E"}
		i := 2

		// Remove the element at index i from slice.
		slice[i] = slice[len(slice)-1] // Copy last element to index i.
		slice[len(slice)-1] = ""       // Erase last element (write zero value).
		slice = slice[:len(slice)-1]   // Truncate slice.
	}
}

func BenchmarkRemoveElementWithCopy(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		slice := []string{"A", "B", "C", "D", "E"}
		i := 2

		// Remove the element at index i from slice.
		copy(slice[i:], slice[i+1:]) // Shift slice[i+1:] left one index.
		slice[len(slice)-1] = ""     // Erase last element (write zero value).
		slice = slice[:len(slice)-1] // Truncate slice.
	}
}

func BenchmarkSimpleRemoveElement(b *testing.B) {
	for iteration := 0; iteration < b.N; iteration++ {
		slice := []string{"A", "B", "C", "D", "E"}
		i := 2
		slice = append(slice[:i], slice[i+1:]...)
	}
}
