package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"net"
)

const MXValidatorName = "MXValidatorInterface"

type MXs = []*net.MX

type MXValidationResultInterface interface {
	MX() MXs
	ValidationResultInterface
}

type MXValidationResult struct {
	mx MXs
	AValidationResult
}

func (v MXValidationResult) MX() MXs {
	return v.mx
}

type MXValidatorInterface interface {
	ValidatorInterface
}

type MXValidator struct{}

func (v MXValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	var mxs MXs
	var err error
	mxs, err = net.LookupMX(email.GetDomain())

	return MXValidationResult{
		mxs,
		NewValidatorResult(err == nil, []error{err}, nil).(ValidationResult),
	}
}
