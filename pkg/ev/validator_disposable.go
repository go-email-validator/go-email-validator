package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

const DisposableValidatorName ValidatorName = "DisposableValidator"

func NewDisposableValidator(d disposable.Interface) ValidatorInterface {
	return DisposableValidator{d: d}
}

type DisposableValidator struct {
	d disposable.Interface
	AValidatorWithoutDeps
}

func (d DisposableValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	return NewValidatorResult(d.d.Disposable(email), nil, nil, DisposableValidatorName)
}
