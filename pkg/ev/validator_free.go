package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
)

const FreeValidatorName ValidatorName = "FreeValidator"

func NewFreeValidator(f free.Interface) ValidatorInterface {
	return FreeValidator{f: f}
}

type FreeValidator struct {
	f free.Interface
	AValidatorWithoutDeps
}

func (r FreeValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	return NewValidatorResult(r.f.IsFree(email), nil, nil, FreeValidatorName)
}
