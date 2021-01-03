package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net/mail"
	"regexp"
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

var defaultEmailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

// Create SyntaxValidatorName, based on *regexp.Regexp
// Example of regular expressions
// HTML5 - https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
// Different variations - http://emailregex.com/, RFC 5322 is used as default emailRegex
func NewSyntaxRegexValidator(emailRegex *regexp.Regexp) Validator {
	if emailRegex == nil {
		emailRegex = defaultEmailRegex
	}

	return syntaxRegexValidator{
		emailRegex: emailRegex,
	}
}

type syntaxRegexValidator struct {
	AValidatorWithoutDeps
	emailRegex *regexp.Regexp
}

func (s syntaxRegexValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	if s.emailRegex.MatchString(email.String()) {
		return NewValidResult(SyntaxValidatorName)
	}
	return syntaxGetError()
}
