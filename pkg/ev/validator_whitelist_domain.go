package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const WhiteListDomainValidatorName ValidatorName = "WhiteListDomain"

type WhiteListError struct {
	utils.Err
}

type WhiteListValidator struct {
	d contains.Interface
	AValidatorWithoutDeps
}

func (w WhiteListValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isContains = w.d.Contains(email.Domain())
	if isContains {
		err = &WhiteListError{}
	}

	return NewValidatorResult(isContains, utils.Errs(err), nil, WhiteListDomainValidatorName)
}
