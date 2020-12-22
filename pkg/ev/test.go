package ev

import "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"

const (
	validUsername = "go.email.validator"
	validDomain   = "gmail.com"
)

func GetValidTestEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress(validUsername, validDomain)
}
