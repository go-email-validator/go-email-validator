package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"net"
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

func NewMXValidator() Validator {
	return mxValidator{}
}

type mxValidator struct{ AValidatorWithoutDeps }

func (v mxValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var mxs evsmtp.MXs
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
