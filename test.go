// Golang program to illustrate
// reflect.Copy() Function

package main

import (
	"fmt"
)

type name struct {
}

func retName() *name {
	return nil
}

func ptrString() interface{} {
	return retName()
}

func String() interface{} {
	return nil
}

// Main function
func main() {
	n := interface{}(nil)

	ptr := ptrString()
	str := String()

	fmt.Println(ptr == nil, ptr == n, str == nil, str == n)
}
