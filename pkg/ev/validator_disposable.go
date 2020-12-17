package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

type DisposableError struct {
	error
}

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return DisposableValidator{d: d}
}

type DisposableValidator struct {
	d disposable.Interface
	AValidatorWithoutDeps
}

func (d DisposableValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var errs = make([]error, 0)
	var isDisposable = d.d.Disposable(email)
	if isDisposable {
		errs = append(errs, DisposableError{})
	}

	return NewValidatorResult(!isDisposable, errs, nil, DisposableValidatorName)
}
