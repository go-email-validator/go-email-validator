package domain

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	Contains(email ev_email.EmailAddressInterface) bool
}

type SetDomain struct {
	set sets.Set
}

func (s SetDomain) Contains(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Domain())
}
