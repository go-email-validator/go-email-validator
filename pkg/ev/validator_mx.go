package ev

import "net"

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

type MXValidator struct{}

func (v MXValidator) Validate(email EmailAddressInterface) ValidationResultInterface {
	var mxs MXs
	var err error
	mxs, err = net.LookupMX(email.GetDomain())

	return MXValidationResult{
		mxs,
		NewValidatorResult(err != nil, []error{err}, nil).(ValidationResult),
	}
}
