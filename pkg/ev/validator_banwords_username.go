package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// BanWordsUsernameValidatorName is name of ban words username
// If email username has ban worlds, the email is invalid
const BanWordsUsernameValidatorName ValidatorName = "BanWordsUsername"

// BanWordsUsernameError is BanWordsUsernameValidatorName error
type BanWordsUsernameError struct{}

func (BanWordsUsernameError) Error() string {
	return "BanWordsUsernameError"
}

// NewBanWordsUsername instantiates BanWordsUsernameValidatorName validator
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
