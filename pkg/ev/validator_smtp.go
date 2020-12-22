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
	mxResult := results[1].(MXValidationResultInterface)

	if syntaxResult.IsValid() && mxResult.IsValid() {
		err := s.checker.Validate(mxResult.MX(), email)

		if err != nil {
			return NewValidatorResult(false, err, nil, SMTPValidatorName)
		}
	}

	return NewValidatorResult(true, nil, nil, SMTPValidatorName)
}
