package ev_email

import (
	"fmt"
	"strings"
)

const (
	AT = "@"
)

type EmailAddress interface {
	Username() string
	Domain() string
	fmt.Stringer
}

func NewEmailAddress(username, domain string) EmailAddress {
	return NewEmailAddressWithSource(username, domain, username+AT+domain)
}

func NewEmailAddressWithSource(username, domain, source string) EmailAddress {
	username = strings.ToLower(username)
	domain = strings.ToLower(domain)
	source = strings.ToLower(source)
	return emailAddress{
		username: username,
		domain:   domain,
		source:   source,
	}
}

type emailAddress struct {
	username string
	domain   string
	source   string
}

func (e emailAddress) Username() string {
	return e.username
}

func (e emailAddress) Domain() string {
	return e.domain
}

func (e emailAddress) String() string {
	return e.source
}

func SeparatedEmail(email string) (string, string) {
	pos := strings.IndexByte(email, '@')

	if pos == -1 || len(email) < 3 {
		return "", ""
	}

	return email[:pos], email[pos+1:]
}

func EmailFromString(email string) EmailAddress {
	username, domain := SeparatedEmail(email)
	return NewEmailAddressWithSource(username, domain, email)
}
