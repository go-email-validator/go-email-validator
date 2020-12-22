package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BlackListDomainsValidatorName ValidatorName = "BlackListDomains"

type BlackListDomainsError struct {
	utils.Err
}

func NewBlackListValidator(d contains.Interface) Validator {
	return blackListValidator{d: d}
}

type blackListValidator struct {
	d contains.Interface
	AValidatorWithoutDeps
}

func (w blackListValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.Domain())
	if isContains {
		err = BlackListDomainsError{}
	}

	return NewValidatorResult(!isContains, utils.Errs(err), nil, BlackListDomainsValidatorName)
}
