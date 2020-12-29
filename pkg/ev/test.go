package ev

import "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"

const (
	validUsername    = "go.email.validator"
	validDomain      = "gmail.com"
	validEmailString = validUsername + "@" + validDomain
)

func GetValidTestEmail() evmail.Address {
	return evmail.NewEmailAddress(validUsername, validDomain)
}
