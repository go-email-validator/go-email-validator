package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// BlackListDomainsValidatorName is name of black domain
// It checks an domain of email in list. If the address domain is in, the email is invalid.
const BlackListDomainsValidatorName ValidatorName = "BlackListDomains"

// BlackListDomainsError is BlackListEmailsValidatorName error
type BlackListDomainsError struct{}

func (BlackListDomainsError) Error() string {
	return "BlackListDomainsError"
}

// NewBlackListValidator instantiates BlackListDomainsValidatorName validator
func NewBlackListValidator(d contains.InSet) Validator {
	return blackListValidator{d: d}
}

type blackListValidator struct {
	d contains.InSet
	AValidatorWithoutDeps
}

func (w blackListValidator) Validate(input Input, _ ...ValidationResult) ValidationResult {
	var err error
	var isContains = w.d.Contains(input.Email().Domain())
	if isContains {
		err = BlackListDomainsError{}
	}

	return NewResult(!isContains, utils.Errs(err), nil, BlackListDomainsValidatorName)
}
