package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const FreeValidatorName ValidatorName = "FreeValidator"

type FreeError struct {
	utils.Err
}

func FreeDefaultValidator() Validator {
	return NewFreeValidator(free.NewWillWhiteSetFree())
}

func NewFreeValidator(f contains.InSet) Validator {
	return freeValidator{f: f}
}

type freeValidator struct {
	AValidatorWithoutDeps
	f contains.InSet
}

func (r freeValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isFree = r.f.Contains(email.Domain())
	if isFree {
		err = FreeError{}
	}

	return NewValidatorResult(!isFree, utils.Errs(err), nil, FreeValidatorName)
}
