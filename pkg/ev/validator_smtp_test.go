package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/smtp_checker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newSMTPValidator() *SMTPValidator {

	return &SMTPValidator{
		smtp_checker.Checker{
			GetConn:   smtp_checker.SimpleClientGetter,
			SendMail:  smtp_checker.NewSendMail(),
			FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
		},
		ADepValidator{},
	}
}

func getSmtpValidator_Validate() DepValidator {
	return DepValidator{
		map[string]ValidatorInterface{
			SyntaxValidatorName: &SyntaxValidator{},
			MXValidatorName:     &MXValidator{},
			SMTPValidatorName:   ValidatorInterface(newSMTPValidator()),
		},
	}
}

func BenchmarkSMTPValidator_Validate(b *testing.B) {
	email := ev_email.NewEmail("go.email.validator", "gmail.com")
	depValidator := getSmtpValidator_Validate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}

func TestSMTPValidator_Validate_WithoutMock(t *testing.T) {
	email := ev_email.NewEmail("go.email.validator", "gmail.com")
	depValidator := getSmtpValidator_Validate()

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
