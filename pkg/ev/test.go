package ev

import "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"

func getValidEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("some", "email.valid")
}
