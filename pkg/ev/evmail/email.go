package evmail

import (
	"fmt"
	"strings"
)

const (
	// AT store at symbol
	AT = "@"
)

// Address represents email
type Address interface {
	Username() string
	Domain() string
	fmt.Stringer
}

// NewEmailAddress forms Address from username and domain
func NewEmailAddress(username, domain string) Address {
	return NewEmailAddressWithSource(username, domain, username+AT+domain)
}

// NewEmailAddressWithSource forms Address
// source used to store empty emails or without username or domain part
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

// SeparateEmail separates email by "@" and returns two parts
func SeparateEmail(email string) (string, string) {
	pos := strings.IndexByte(email, '@')

	if pos == -1 || len(email) < 3 {
		return "", ""
	}

	return email[:pos], email[pos+1:]
}

// FromString forms Address from string
func FromString(email string) Address {
	username, domain := SeparateEmail(email)

	return NewEmailAddressWithSource(username, domain, email)
}
