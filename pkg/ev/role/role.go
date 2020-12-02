package role

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
)

type Interface interface {
	HasRole(email ev_email.EmailAddressInterface) bool
}

type SetRole struct {
	set utils.StringSet
}

func (s SetRole) HasRole(email ev_email.EmailAddressInterface) bool {
	_, ok := s.set[email.Username()]
	return ok
}
