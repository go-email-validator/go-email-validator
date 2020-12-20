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

func NewEmailAddress(username, domain string) EmailAddressInterface {
	return emailAddress{
		strings.ToLower(username),
		strings.ToLower(domain),
	}
}

type emailAddress struct {
	username string
	domain   string
}

func (e emailAddress) Username() string {
	return e.username
}

func (e emailAddress) Domain() string {
	return e.domain
}

func (e emailAddress) String() string {
	return e.Username() + AT + e.Domain()
}

func SeparatedEmail(email string) (string, string) {
	pos := strings.Index(email, "@")
	return email[:pos], email[pos+1:]
}

func EmailFromString(email string) EmailAddressInterface {
	return NewEmailAddress(SeparatedEmail(email))
}
