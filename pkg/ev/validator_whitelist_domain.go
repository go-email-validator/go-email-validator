package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const WhiteListDomainValidatorName ValidatorName = "WhiteListDomain"

type WhiteListError struct {
	utils.Err
}

func NewWhiteListValidator(d contains.InSet) Validator {
	return whiteListValidator{d: d}
}

type whiteListValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w whiteListValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.Domain())
	if isContains {
		err = &WhiteListError{}
	}

	return NewValidatorResult(isContains, utils.Errs(err), nil, WhiteListDomainValidatorName)
}
