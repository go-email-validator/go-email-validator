package contains

import (
	"github.com/emirpasic/gods/sets"
)

// InSet checks presence of value
type InSet interface {
	Contains(value interface{}) bool
}

// NewSet instantiates InSet based on sets.Set
func NewSet(s sets.Set) InSet {
	return setContains{s}
}

type setContains struct {
	set sets.Set
}

func (s setContains) Contains(value interface{}) bool {
	return s.set.Contains(value)
}

// FuncChecker checks presence of value
type FuncChecker func(value interface{}) bool

// NewFunc instantiates InSet, functions is used for checking
func NewFunc(f FuncChecker) InSet {
	return funcContains{f}
}

type funcContains struct {
	funcChecker FuncChecker
}

func (f funcContains) Contains(value interface{}) bool {
	return f.funcChecker(value)
}
