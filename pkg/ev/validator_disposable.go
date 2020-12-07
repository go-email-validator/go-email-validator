package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/disposable"
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

const DisposableValidatorName = "DisposableValidator"

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return DisposableValidator{d: d}
}

type DisposableValidator struct {
	d disposable.Interface
	AValidatorWithoutDeps
}

func (d DisposableValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	return NewValidatorResult(d.d.Disposable(email), nil, nil)
}
