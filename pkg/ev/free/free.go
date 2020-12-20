package free

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	IsFree(email ev_email.EmailAddressInterface) bool
}

func NewSetFree(set sets.Set) Interface {
	return setFree{set}
}

type setFree struct {
	set sets.Set
}

func (s setFree) IsFree(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Domain())
}
