package ev

import (
	"net/mail"
)

type SyntaxValidator struct{}

func (s SyntaxValidator) Validate(email EmailAddressInterface) ValidationResultInterface {
	var err error
	_, err = mail.ParseAddress(email.String())

	return NewValidatorResult(err != nil, []error{err}, nil)
}
