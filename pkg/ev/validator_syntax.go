package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net/mail"
)

const SyntaxValidatorName ValidatorName = "syntaxValidator"

type SyntaxError struct {
	error
}

type SyntaxValidatorResultInterface interface {
	ValidationResultInterface
}

func NewSyntaxValidator() ValidatorInterface {
	return syntaxValidator{}
}

type syntaxValidator struct {
	AValidatorWithoutDeps
}

func (_ syntaxValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	_, err := mail.ParseAddress(email.String())

	return NewValidatorResult(err == nil, utils.Errs(&SyntaxError{err}), nil, SyntaxValidatorName)
}
