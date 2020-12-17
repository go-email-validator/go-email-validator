package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
)

const FreeValidatorName ValidatorName = "FreeValidator"

type FreeError struct {
	error
}

func FreeDefaultValidator() ValidatorInterface {
	return NewFreeValidator(free.NewWillWhiteSetFree())
}

func NewFreeValidator(f free.Interface) ValidatorInterface {
	return FreeValidator{f: f}
}

type FreeValidator struct {
	f free.Interface
	AValidatorWithoutDeps
}

func (r FreeValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var errs = make([]error, 0)
	var isFree = r.f.IsFree(email)
	if isFree {
		errs = append(errs, FreeError{})
	}

	return NewValidatorResult(!isFree, errs, nil, FreeValidatorName)
}
