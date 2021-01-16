package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// MXValidatorName  is name of mx validator
const MXValidatorName ValidatorName = "MXValidator"

// EmptyMXsError is error of MXValidatorName
type EmptyMXsError struct{}

func (EmptyMXsError) Error() string {
	return "EmptyMXsError"
}

// MXValidationResult is result of MXValidatorName
type MXValidationResult interface {
	MX() evsmtp.MXs
	ValidationResult
}

// NewMXValidationResult instantiates result of MXValidatorName
func NewMXValidationResult(mx evsmtp.MXs, result *AValidationResult) MXValidationResult {
	return mxValidationResult{mx: mx, AValidationResult: result}
}

type mxValidationResult struct {
	*AValidationResult
	mx evsmtp.MXs
}

func (v mxValidationResult) MX() evsmtp.MXs {
	return v.mx
}

// DefaultNewMXValidator instantiates default MXValidatorName based on evsmtp.LookupMX
func DefaultNewMXValidator() Validator {
	return NewMXValidator(evsmtp.LookupMX)
}

// NewMXValidator instantiates MXValidatorName based on evsmtp.FuncLookupMX
func NewMXValidator(lookupMX evsmtp.FuncLookupMX) Validator {
	return mxValidator{
		lookupMX: lookupMX,
	}
}

type mxValidator struct {
	AValidatorWithoutDeps
	lookupMX evsmtp.FuncLookupMX
}

func (v mxValidator) Validate(input Interface, _ ...ValidationResult) ValidationResult {
	var mxs evsmtp.MXs
	var err error
	mxs, err = v.lookupMX(input.Email().Domain())

	if hasMXs := len(mxs) > 0; err == nil && !hasMXs {
		err = EmptyMXsError{}
	}

	return NewMXValidationResult(
		mxs,
		NewResult(err == nil, utils.Errs(err), nil, MXValidatorName).(*AValidationResult),
	)
}
