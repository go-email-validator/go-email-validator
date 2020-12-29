package utils

import (
	"go/types"
)

func RangeLen(i interface{}) int {
	if i == nil {
		return 0
	}

	switch i.(type) {
	case types.Array:
	case types.Slice:
		return len(i.([]interface{}))
	case types.Map:
		return len(i.(map[interface{}]interface{}))
	}
	panic("interface should be \"array\", \"slice\" or \"map\"")
}

func Errs(errs ...error) []error {
	if errs == nil || len(errs) == 1 && errs[0] == nil {
		return nil
	}

	return errs
}

func NewError(text string) error {
	return Err{text}
}

type Err struct {
	s string
}

func (e Err) Error() string {
	return e.s
}
