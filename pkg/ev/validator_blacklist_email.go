package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BlackListEmailsValidatorName ValidatorName = "BlackListEmails"

type BlackListEmailsError struct {
	utils.Err
}

func NewBlackListEmailsValidator(d contains.InSet) Validator {
	return blackListEmailsValidator{d: d}
}

type blackListEmailsValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w blackListEmailsValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.String())
	if isContains {
		err = BlackListEmailsError{}
	}

	return NewValidatorResult(!isContains, utils.Errs(err), nil, BlackListEmailsValidatorName)
}
