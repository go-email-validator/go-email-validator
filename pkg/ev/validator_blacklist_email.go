package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// BlackListEmailsValidatorName is name of black list emails
// It checks an email in list. If the address is in, the email is invalid.
const BlackListEmailsValidatorName ValidatorName = "BlackListEmails"

// BlackListEmailsErr is text for BlackListEmailsError.Error
const BlackListEmailsErr = "BlackListEmailsError"

// BlackListEmailsError is BlackListEmailsValidatorName error
type BlackListEmailsError struct{}

func (BlackListEmailsError) Error() string {
	return BlackListEmailsErr
}

// NewBlackListEmailsValidator instantiates BlackListEmailsValidatorName validator
func NewBlackListEmailsValidator(d contains.InSet) Validator {
	return blackListEmailsValidator{d: d}
}

type blackListEmailsValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w blackListEmailsValidator) Validate(input Input, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(input.Email().String())
	if isContains {
		err = BlackListEmailsError{}
	}

	return NewResult(!isContains, utils.Errs(err), nil, BlackListEmailsValidatorName)
}
