package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/domain"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const WhiteListValidatorName ValidatorName = "WhiteList"

type WhiteListError struct {
	utils.Err
}

type WhiteListValidator struct {
	d domain.Interface
	AValidatorWithoutDeps
}

func (w WhiteListValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isContains = w.d.Contains(email)
	if isContains {
		err = &WhiteListError{}
	}

	return NewValidatorResult(isContains, utils.Errs(err), nil, WhiteListValidatorName)
}
