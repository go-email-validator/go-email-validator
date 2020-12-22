package contains

import (
	"github.com/emirpasic/gods/sets"
)

type Interface interface {
	Contains(value interface{}) bool
}

func NewSet(s sets.Set) Interface {
	return setContains{s}
}

type setContains struct {
	set sets.Set
}

func (s setContains) Contains(value interface{}) bool {
	return s.set.Contains(value)
}

type FuncChecker func(value interface{}) bool

func NewFunc(f FuncChecker) Interface {
	return funcContains{f}
}

type funcContains struct {
	funcChecker FuncChecker
}

func (f funcContains) Contains(value interface{}) bool {
	return f.funcChecker(value)
}
