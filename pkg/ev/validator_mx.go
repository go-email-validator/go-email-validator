package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const MXValidatorName ValidatorName = "MXValidator"

type EmptyMXsError struct {
	utils.Err
}

type MXValidationResult interface {
	MX() evsmtp.MXs
	ValidationResult
}

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

func DefaultNewMXValidator() Validator {
	return NewMXValidator(evsmtp.LookupMX)
}

func NewMXValidator(lookupMX evsmtp.FuncLookupMX) Validator {
	return mxValidator{
		lookupMX: lookupMX,
	}
}

type mxValidator struct {
	AValidatorWithoutDeps
	lookupMX evsmtp.FuncLookupMX
}

func (v mxValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var mxs evsmtp.MXs
	var err error
	mxs, err = v.lookupMX(email.Domain())

	if hasMXs := len(mxs) > 0; err == nil && !hasMXs {
		err = EmptyMXsError{}
	}

	return NewMXValidationResult(
		mxs,
		NewResult(err == nil, utils.Errs(err), nil, MXValidatorName).(*AValidationResult),
	)
}
