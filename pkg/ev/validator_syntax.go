package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"net/mail"
)

const SyntaxValidatorName ValidatorName = "SyntaxValidator"

type SyntaxValidatorResultInterface interface {
	ValidationResultInterface
}

type SyntaxValidator struct {
	AValidatorWithoutDeps
}

func (s SyntaxValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	_, err = mail.ParseAddress(email.String())

	return NewValidatorResult(err == nil, []error{err}, nil, SyntaxValidatorName)
}
