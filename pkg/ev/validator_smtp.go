package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
)

const SMTPValidatorName ValidatorName = "SMTPValidator"

func NewSMTPValidator(Checker smtp_checker.CheckerInterface) Validator {
	return smtpValidator{Checker}
}

type smtpValidator struct {
	checker smtp_checker.CheckerInterface
}

func (s smtpValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName, MXValidatorName}
}

func (s smtpValidator) Validate(email ev_email.EmailAddress, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResultInterface)
	mxResult := results[1].(MXValidationResult)
	var errs []error

	isDepsValid := syntaxResult.IsValid() && mxResult.IsValid()
	if isDepsValid {
		errs = s.checker.Validate(mxResult.MX(), email)
	}

	return NewValidatorResult(isDepsValid && len(errs) == 0, errs, nil, SMTPValidatorName)
}
