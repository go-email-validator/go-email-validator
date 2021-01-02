package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net/mail"
)

const SyntaxValidatorName ValidatorName = "syntaxValidator"

type SyntaxError struct {
	error
}

type SyntaxValidatorResult interface {
	ValidationResult
}

func NewSyntaxValidator() Validator {
	return syntaxValidator{}
}

type syntaxValidator struct {
	AValidatorWithoutDeps
}

func (_ syntaxValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	_, err := mail.ParseAddress(email.String())

	if err == nil {
		return NewValidResult(SyntaxValidatorName)
	}
	return syntaxGetError()
}

func syntaxGetError() ValidationResult {
	return NewResult(false, utils.Errs(SyntaxError{}), nil, SyntaxValidatorName)
}
