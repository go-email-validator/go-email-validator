package evmail

import (
	"fmt"
	"strings"
)

const (
	AT = "@"
)

type Address interface {
	Username() string
	Domain() string
	fmt.Stringer
}

func NewEmailAddress(username, domain string) Address {
	return NewEmailAddressWithSource(username, domain, username+AT+domain)
}

func NewEmailAddressWithSource(username, domain, source string) Address {
	username = strings.ToLower(username)
	domain = strings.ToLower(domain)
	source = strings.ToLower(source)

	return address{
		username: username,
		domain:   domain,
		source:   source,
	}
}

type address struct {
	username string
	domain   string
	source   string
}

func (e address) Username() string {
	return e.username
}

func (e address) Domain() string {
	return e.domain
}

func (e address) String() string {
	return e.source
}

func SeparateEmail(email string) (string, string) {
	pos := strings.IndexByte(email, '@')

	if pos == -1 || len(email) < 3 {
		return "", ""
	}

	return email[:pos], email[pos+1:]
}

func FromString(email string) Address {
	username, domain := SeparateEmail(email)

	return NewEmailAddressWithSource(username, domain, email)
}
