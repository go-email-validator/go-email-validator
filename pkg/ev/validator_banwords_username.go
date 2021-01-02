package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const BanWordsUsernameValidatorName ValidatorName = "BanWordsUsername"

type BanWordsUsernameError struct {
	utils.Err
}

func NewBanWordsUsername(inStrings contains.InStrings) Validator {
	return banWordsUsernameValidator{d: inStrings}
}

type banWordsUsernameValidator struct {
	d contains.InStrings
	AValidatorWithoutDeps
}

func (w banWordsUsernameValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(email.Username())
	if isContains {
		err = BanWordsUsernameError{}
	}

	return NewResult(!isContains, utils.Errs(err), nil, BanWordsUsernameValidatorName)
}
