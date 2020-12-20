package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const FreeValidatorName ValidatorName = "FreeValidator"

type FreeError struct {
	utils.Err
}

func FreeDefaultValidator() ValidatorInterface {
	return NewFreeValidator(free.NewWillWhiteSetFree())
}

func NewFreeValidator(f free.Interface) ValidatorInterface {
	return freeValidator{f: f}
}

type freeValidator struct {
	f free.Interface
	AValidatorWithoutDeps
}

func (r freeValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isFree = r.f.IsFree(email)
	if isFree {
		err = FreeError{}
	}

	return NewValidatorResult(!isFree, utils.Errs(err), nil, FreeValidatorName)
}
