package ev

import "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"

func getValidEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmail("some", "email.valid")
}
