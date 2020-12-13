package free

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	IsFree(email ev_email.EmailAddressInterface) bool
}

type SetFree struct {
	set sets.Set
}

func (s SetFree) IsFree(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Domain())
}
