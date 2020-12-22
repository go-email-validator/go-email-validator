package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

type DisposableError struct {
	utils.Err
}

func NewDisposableValidator(d contains.Interface) Validator {
	return disposableValidator{d: d}
}

type disposableValidator struct {
	d contains.Interface
	AValidatorWithoutDeps
}

func (d disposableValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	var isDisposable = d.d.Contains(email.Domain())
	if isDisposable {
		err = DisposableError{}
	}

	return NewValidatorResult(!isDisposable, utils.Errs(err), nil, DisposableValidatorName)
}
