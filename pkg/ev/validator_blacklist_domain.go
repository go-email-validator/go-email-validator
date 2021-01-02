package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BlackListDomainsValidatorName ValidatorName = "BlackListDomains"

type BlackListDomainsError struct {
	utils.Err
}

func NewBlackListValidator(d contains.InSet) Validator {
	return blackListValidator{d: d}
}

type blackListValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w blackListValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.Domain())
	if isContains {
		err = BlackListDomainsError{}
	}

	return NewResult(!isContains, utils.Errs(err), nil, BlackListDomainsValidatorName)
}
