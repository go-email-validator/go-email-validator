package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
)

const SMTPValidatorName ValidatorName = "SMTPValidator"

func NewSMTPValidator(Checker smtp_checker.CheckerInterface) ValidatorInterface {
	return SMTPValidator{Checker}
}

type SMTPValidator struct {
	Checker smtp_checker.CheckerInterface
}

func (a SMTPValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName, MXValidatorName}
}

func (a SMTPValidator) Validate(email ev_email.EmailAddressInterface, results ...ValidationResultInterface) ValidationResultInterface {
	syntaxResult := results[0].(SyntaxValidatorResultInterface)
	mxResult := results[1].(MXValidationResultInterface)

	if syntaxResult.IsValid() && mxResult.IsValid() {
		err := a.Checker.Validate(mxResult.MX(), email)

		if err != nil {
			return NewValidatorResult(false, err, nil, SMTPValidatorName)
		}
	}

	return NewValidatorResult(true, nil, nil, SMTPValidatorName)
}
