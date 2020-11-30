package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/disposable"
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

const DisposableValidatorName = "DisposableValidatorInterface"

type DisposableValidatorInterface interface {
	ValidatorInterface
}

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return DisposableValidator{d}
}

type DisposableValidator struct {
	d disposable.Interface
}

func (d DisposableValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	return NewValidatorResult(d.d.Disposable(email), nil, nil)
}
