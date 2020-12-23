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
	ValidationResult
}

func NewSyntaxValidator() Validator {
	return syntaxValidator{}
}

type syntaxValidator struct {
	AValidatorWithoutDeps
}

func (_ syntaxValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	_, err := mail.ParseAddress(email.String())

	if err == nil {
		return NewValidValidatorResult(SyntaxValidatorName)
	}
	return syntaxGetError()
}

func syntaxGetError() ValidationResult {
	return NewValidatorResult(false, utils.Errs(SyntaxError{}), nil, SyntaxValidatorName)
}
