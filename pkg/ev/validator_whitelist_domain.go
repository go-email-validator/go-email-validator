package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const WhiteListDomainValidatorName ValidatorName = "WhiteListDomain"

type WhiteListError struct {
	utils.Err
}

func NewWhiteListValidator(d contains.Interface) Validator {
	return whiteListValidator{d: d}
}

type whiteListValidator struct {
	d contains.Interface
	AValidatorWithoutDeps
}

func (w whiteListValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.Domain())
	if isContains {
		err = &WhiteListError{}
	}

	return NewValidatorResult(isContains, utils.Errs(err), nil, WhiteListDomainValidatorName)
}
