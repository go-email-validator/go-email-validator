package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
)

const SMTPValidatorName ValidatorName = "SMTPValidator"

func NewSMTPValidator(Checker evsmtp.Checker) Validator {
	return smtpValidator{Checker}
}

type smtpValidator struct {
	checker evsmtp.Checker
}

func (s smtpValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName, MXValidatorName}
}

func (s smtpValidator) Validate(email evmail.Address, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResult)
	mxResult := results[1].(MXValidationResult)
	var errs []error

	if syntaxResult.IsValid() && mxResult.IsValid() {
		errs = s.checker.Validate(mxResult.MX(), email)
	} else {
		errs = append(errs, DepsError{})
	}

	return NewValidatorResult(len(errs) == 0, errs, nil, SMTPValidatorName)
}
