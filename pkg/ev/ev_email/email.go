package ev_email

import (
	"fmt"
	"strings"
)

type EmailAddressInterface interface {
	GetUsername() string
	GetDomain() string
	fmt.Stringer
}

type EmailAddress struct {
	username string
	domain   string
	atSymbol string
}

func (e EmailAddress) GetUsername() string {
	return e.username
}

func (e EmailAddress) GetDomain() string {
	return e.domain
}

func (e EmailAddress) String() string {
	return e.GetUsername() + e.atSymbol + e.GetDomain()
}

func NewEmailAddress(username, domain string) EmailAddressInterface {
	return EmailAddress{
		strings.ToLower(username),
		strings.ToLower(domain),
		"@",
	}
}
