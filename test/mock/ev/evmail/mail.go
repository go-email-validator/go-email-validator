package mockevmail

import "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"

const (
	ValidUsername    = "go.email.validator"
	ValidDomain      = "gmail.com"
	ValidEmailString = ValidUsername + "@" + ValidDomain
)

// GetValidTestEmail returns valid email.Address
func GetValidTestEmail() evmail.Address {
	return evmail.NewEmailAddress(ValidUsername, ValidDomain)
}
