package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net/mail"
	"regexp"
)

// SyntaxValidatorName is name of syntax validator
const SyntaxValidatorName ValidatorName = "syntaxValidator"

// SyntaxError is error of SyntaxValidatorName
type SyntaxError struct{}

func (SyntaxError) Error() string {
	return "SyntaxError"
}

// SyntaxValidatorResult is interface of SyntaxValidatorName result
type SyntaxValidatorResult interface {
	ValidationResult
}

// NewSyntaxValidator instantiates SyntaxValidatorName based on mail.ParseAddress
func NewSyntaxValidator() Validator {
	return syntaxValidator{}
}

type syntaxValidator struct {
	AValidatorWithoutDeps
}

func (syntaxValidator) Validate(input Interface, _ ...ValidationResult) ValidationResult {
	_, err := mail.ParseAddress(input.Email().String())

	if err == nil {
		return NewValidResult(SyntaxValidatorName)
	}
	return syntaxGetError()
}

func syntaxGetError() ValidationResult {
	return NewResult(false, utils.Errs(SyntaxError{}), nil, SyntaxValidatorName)
}

var defaultEmailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")

// NewSyntaxRegexValidator creates SyntaxValidatorName, based on *regexp.Regexp
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

func (s syntaxRegexValidator) Validate(input Interface, _ ...ValidationResult) ValidationResult {
	if s.emailRegex.MatchString(input.Email().String()) {
		return NewValidResult(SyntaxValidatorName)
	}
	return syntaxGetError()
}
