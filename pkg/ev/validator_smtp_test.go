package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
	"github.com/stretchr/testify/assert"
	"net/smtp"
	"testing"
)

//test monicaramirezrestrepo@hotmail.com
func newSMTPValidator() *smtpValidator {
	return &smtpValidator{
		checker: smtp_checker.NewChecker(smtp_checker.CheckerDTO{
			DialFunc:  smtp.Dial,
			SendMail:  smtp_checker.NewSendMail(),
			FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
		}),
	}
}

func getSmtpValidator_Validate() Validator {
	return NewDepValidator(
		map[ValidatorName]Validator{
			SyntaxValidatorName: NewSyntaxValidator(),
			MXValidatorName:     NewMXValidator(),
			SMTPValidatorName: NewWarningsDecorator(
				Validator(newSMTPValidator()),
				NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
					return func(err error) bool {
						return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
					}
				}),
			),
		},
	)
}

func BenchmarkSMTPValidator_Validate(b *testing.B) {
	email := ev_email.EmailFromString(validEmailString)
	depValidator := getSmtpValidator_Validate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}

func TestSMTPValidator_Validate_WithoutMock(t *testing.T) {
	email := ev_email.EmailFromString(validEmailString)
	depValidator := getSmtpValidator_Validate()

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
