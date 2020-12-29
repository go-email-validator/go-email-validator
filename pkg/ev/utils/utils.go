package utils

import (
	"reflect"
)

func RangeLen(i interface{}) int {
	if i == nil {
		return 0
	}

	return reflect.ValueOf(i).Len()
}

func Errs(errs ...error) []error {
	if errs == nil || len(errs) == 1 && errs[0] == nil {
		return nil
	}

	return errs
}

type Err struct {
	s string
}

func (e Err) Error() string {
	return e.s
}
