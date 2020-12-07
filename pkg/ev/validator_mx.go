package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
	"net"
)

const MXValidatorName = "MXValidator"

type MXValidationResultInterface interface {
	MX() utils.MXs
	ValidationResultInterface
}

type MXValidationResult struct {
	mx utils.MXs
	*AValidationResult
}

func (v MXValidationResult) MX() utils.MXs {
	return v.mx
}

type MXValidator struct{ AValidatorWithoutDeps }

func (v MXValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var mxs utils.MXs
	var err error
	mxs, err = net.LookupMX(email.Domain())

	// TODO fix []error{err}
	return MXValidationResult{
		mxs,
		NewValidatorResult(err == nil, []error{err}, nil).(*AValidationResult),
	}
}
