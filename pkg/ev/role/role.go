package role

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	HasRole(email ev_email.EmailAddressInterface) bool
}

type SetRole struct {
	set sets.Set
}

func (s SetRole) HasRole(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Username())
}
