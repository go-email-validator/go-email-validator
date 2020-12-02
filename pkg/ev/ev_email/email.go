package ev_email

import (
	"fmt"
	"strings"
)

const (
	AT = "@"
)

type EmailAddressInterface interface {
	Username() string
	Domain() string
	fmt.Stringer
}

type EmailAddress struct {
	username string
	domain   string
}

func (e EmailAddress) Username() string {
	return e.username
}

func (e EmailAddress) Domain() string {
	return e.domain
}

func (e EmailAddress) String() string {
	return e.Username() + AT + e.Domain()
}

func SeparatedEmail(email string) (string, string) {
	pos := strings.Index(email, "@")
	return email[:pos], email[pos+1:]
}

func EmailFromString(email string) EmailAddressInterface {
	return NewEmail(SeparatedEmail(email))
}

func NewEmail(username, domain string) EmailAddressInterface {
	return EmailAddress{
		strings.ToLower(username),
		strings.ToLower(domain),
	}
}
