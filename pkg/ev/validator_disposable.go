package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// DisposableValidatorName is name of disposable validator
const DisposableValidatorName ValidatorName = "DisposableValidator"

// DisposableError is DisposableValidatorName error
type DisposableError struct{}

func (DisposableError) Error() string {
	return "DisposableError"
}

// NewDisposableValidator instantiates DisposableValidatorName
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
