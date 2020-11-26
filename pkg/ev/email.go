package ev

import "fmt"

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
