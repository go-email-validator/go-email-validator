package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
)

// SMTPValidatorName is name of smtp validator
const SMTPValidatorName ValidatorName = "SMTPValidator"

// NewSMTPValidator instantiates SMTPValidatorName
func NewSMTPValidator(Checker evsmtp.Checker) Validator {
	return smtpValidator{Checker}
}

type smtpValidator struct {
	checker evsmtp.Checker
}

func (s smtpValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName, MXValidatorName}
}

func (s smtpValidator) Validate(input Input, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResult)
	mxResult := results[1].(MXValidationResult)
	var errs []error

	if syntaxResult.IsValid() && mxResult.IsValid() {
		var opts evsmtp.Options
		if optsInterface := input.Option(SMTPValidatorName); optsInterface != nil {
			opts = optsInterface.(evsmtp.Options)
		}

		errs = s.checker.Validate(
			mxResult.MX(),
			evsmtp.NewInput(input.Email(), opts),
		)
	} else {
		errs = append(errs, NewDepsError())
	}

	return NewResult(len(errs) == 0, errs, nil, SMTPValidatorName)
}
