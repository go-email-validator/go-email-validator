package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
	"net"
)

const MXValidatorName = "MXValidatorInterface"

type MXValidationResultInterface interface {
	MX() utils.MXs
	ValidationResultInterface
}

type MXValidationResult struct {
	mx utils.MXs
	AValidationResult
}

func (v MXValidationResult) MX() utils.MXs {
	return v.mx
}

type MXValidatorInterface interface {
	ValidatorInterface
}

type MXValidator struct{}

func (v MXValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	var mxs utils.MXs
	var err error
	mxs, err = net.LookupMX(email.Domain())

	return MXValidationResult{
		mxs,
		NewValidatorResult(err == nil, []error{err}, nil).(ValidationResult),
	}
}
