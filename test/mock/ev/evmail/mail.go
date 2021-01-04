package mock_evmail

import "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"

const (
	ValidUsername    = "go.email.validator"
	ValidDomain      = "gmail.com"
	ValidEmailString = ValidUsername + "@" + ValidDomain
)

func GetValidTestEmail() evmail.Address {
	return evmail.NewEmailAddress(ValidUsername, ValidDomain)
}
