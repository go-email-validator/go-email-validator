package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// WhiteListDomainValidatorName name of whiteListValidator
const WhiteListDomainValidatorName ValidatorName = "WhiteListDomain"

// WhiteListErr is text for WhiteListError.Error
const WhiteListErr = "WhiteListError"

// WhiteListError is error for WhiteListDomainValidatorName
type WhiteListError struct{}

func (WhiteListError) Error() string {
	return WhiteListErr
}

// NewWhiteListValidator instantiates WhiteListDomainValidatorName
func NewWhiteListValidator(d contains.InSet) Validator {
	return whiteListValidator{d: d}
}

type whiteListValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w whiteListValidator) Validate(input Input, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(input.Email().Domain())
	if !isContains {
		err = WhiteListError{}
	}

	return NewResult(isContains, utils.Errs(err), nil, WhiteListDomainValidatorName)
}
