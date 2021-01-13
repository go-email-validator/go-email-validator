package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// FreeValidatorName is name of free validator
const FreeValidatorName ValidatorName = "FreeValidator"

// FreeError is FreeValidatorName error
type FreeError struct{}

func (FreeError) Error() string {
	return "FreeError"
}

// FreeDefaultValidator instantiates default FreeValidatorName based on free.NewWillWhiteSetFree()
func FreeDefaultValidator() Validator {
	return NewFreeValidator(free.NewWillWhiteSetFree())
}

// NewFreeValidator instantiates FreeValidatorName
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

	return NewResult(!isFree, utils.Errs(err), nil, FreeValidatorName)
}
