package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net"
)

const MXValidatorName ValidatorName = "MXValidator"

type EmptyMXsError struct {
	utils.Err
}

type MXValidationResult interface {
	MX() utils.MXs
	ValidationResult
}

func NewMXValidationResult(mx utils.MXs, result *AValidationResult) MXValidationResult {
	return mxValidationResult{mx, result}
}

type mxValidationResult struct {
	mx utils.MXs
	*AValidationResult
}

func (v mxValidationResult) MX() utils.MXs {
	return v.mx
}

func NewMXValidator() Validator {
	return mxValidator{}
}

type mxValidator struct{ AValidatorWithoutDeps }

func (v mxValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var mxs utils.MXs
	var err error
	mxs, err = net.LookupMX(email.Domain())

	hasMXs := len(mxs) > 0
	if !hasMXs {
		err = EmptyMXsError{}
	}

	return NewMXValidationResult(
		mxs,
		NewValidatorResult(err == nil, utils.Errs(err), nil, MXValidatorName).(*AValidationResult),
	)
}
