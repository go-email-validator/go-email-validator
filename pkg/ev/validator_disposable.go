package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

type DisposableError struct {
	utils.Error
}

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return DisposableValidator{d: d}
}

type DisposableValidator struct {
	d disposable.Interface
	AValidatorWithoutDeps
}

func (d DisposableValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isDisposable = d.d.Disposable(email)
	if isDisposable {
		err = DisposableError{}
	}

	return NewValidatorResult(!isDisposable, utils.Errs(err), nil, DisposableValidatorName)
}
