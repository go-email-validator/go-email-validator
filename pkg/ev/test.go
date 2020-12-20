package ev

import "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"

func GetValidTestEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("go.email.validator", "gmail.com")
}
