package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/domain"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BlackListValidatorName ValidatorName = "BlackList"

type BlackListError struct {
	utils.Error
}

type BlackListValidator struct {
	d domain.Interface
	AValidatorWithoutDeps
}

func (w BlackListValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var isContains = w.d.Contains(email)
	if isContains {
		err = BlackListError{}
	}

	return NewValidatorResult(!isContains, utils.Errs(err), nil, BlackListValidatorName)
}
