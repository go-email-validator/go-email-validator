package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

type DisposableError struct {
	utils.Err
}

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return disposableValidator{d: d}
}

type disposableValidator struct {
	d disposable.Interface
	AValidatorWithoutDeps
}

func (d disposableValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isDisposable = d.d.Disposable(email)
	if isDisposable {
		err = DisposableError{}
	}

	return NewValidatorResult(!isDisposable, utils.Errs(err), nil, DisposableValidatorName)
}
