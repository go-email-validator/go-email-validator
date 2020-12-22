package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BlackListEmailsValidatorName ValidatorName = "BlackListEmails"

type BlackListEmailsError struct {
	utils.Err
}

func NewBlackListEmailsValidator(d contains.Interface) Validator {
	return blackListValidator{d: d}
}

type blackListEmailsValidator struct {
	d contains.Interface
	AValidatorWithoutDeps
}

func (w blackListEmailsValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.String())
	if isContains {
		err = BlackListEmailsError{}
	}

	return NewValidatorResult(!isContains, utils.Errs(err), nil, BlackListEmailsValidatorName)
}
