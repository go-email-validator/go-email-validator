package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

type DisposableError struct {
	utils.Err
}

func NewDisposableValidator(d contains.InSet) Validator {
	return disposableValidator{d: d}
}

type disposableValidator struct {
	AValidatorWithoutDeps
	d contains.InSet
}

func (d disposableValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isDisposable = d.d.Contains(email.Domain())
	if isDisposable {
		err = DisposableError{}
	}

	return NewResult(!isDisposable, utils.Errs(err), nil, DisposableValidatorName)
}
