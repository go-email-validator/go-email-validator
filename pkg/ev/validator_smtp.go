package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/smtp_checker"
)

const SMTPValidatorName = "SMTPValidator"

type SMTPValidator struct {
	Checker smtp_checker.CheckerInterface
}

func (a SMTPValidator) GetDeps() []string {
	return []string{SyntaxValidatorName, MXValidatorName}
}

func (a SMTPValidator) Validate(email ev_email.EmailAddressInterface, results ...ValidationResultInterface) ValidationResultInterface {
	syntaxResult := results[0].(SyntaxValidatorResultInterface)
	mxResult := results[1].(MXValidationResultInterface)

	if syntaxResult.IsValid() && mxResult.IsValid() {
		err := a.Checker.Validate(mxResult.MX(), email)

		if err != nil {
			return NewValidatorResult(false, []error{err}, nil)
		}
	}

	return NewValidatorResult(true, nil, nil)
}
