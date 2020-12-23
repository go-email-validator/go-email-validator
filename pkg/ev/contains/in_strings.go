package contains

import "github.com/emirpasic/gods/sets/hashset"

type InStrings interface {
	Contains(value string) bool
}

func NewInStringsFromArray(elements []string) InStrings {
	var maxLen = 0
	setElements := make([]interface{}, len(elements))
	for i, element := range elements {
		currentLen := len(element)
		if currentLen > maxLen {
			maxLen = currentLen
		}
		setElements[i] = element
	}

	return NewInStrings(NewSet(hashset.New(setElements...)), maxLen)
}

func NewInStrings(contains Interface, maxLen int) InStrings {
	return inStrings{contains, maxLen}
}

type inStrings struct {
	contains Interface
	maxLen   int
}

func (is inStrings) Contains(value string) bool {
	var maxLen = len(value)
	var jEnd int

	for i := 0; i < maxLen; i++ {
		key := ""
		jEnd = i + is.maxLen
		if jEnd > maxLen {
			jEnd = maxLen
		}
		for j := i; j < jEnd; j++ {
			key = key + string(value[j])
			if is.contains.Contains(key) {
				return true
			}
		}
	}

	return false
}
