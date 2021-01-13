package mockevmail

import "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"

const (
	// ValidUsername ...
	ValidUsername = "go.email.validator"
	// ValidDomain ...
	ValidDomain = "gmail.com"
	// ValidEmailString ...
	ValidEmailString = ValidUsername + "@" + ValidDomain
)

// GetValidTestEmail returns valid email.Address
func GetValidTestEmail() evmail.Address {
	return evmail.NewEmailAddress(ValidUsername, ValidDomain)
}
