package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getSmtpValidator_Validate() DepValidator {
	return DepValidator{
		map[string]ValidatorInterface{
			SyntaxValidatorName: &SyntaxValidator{},
			MXValidatorName:     &MXValidator{},
			SMTPValidatorName: ValidatorInterface(NewSMTPValidator(
				nil,
				nil,
			)),
		},
	}
}

func BenchmarkSMTPValidator_Validate(b *testing.B) {
	email := ev_email.NewEmailAddress("ilia.sergunin", "gmail.com")
	depValidator := getSmtpValidator_Validate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}

func TestSMTPValidator_Validate_WithoutMock(t *testing.T) {
	email := ev_email.NewEmailAddress("ilia.sergunin", "gmail.com")
	depValidator := getSmtpValidator_Validate()

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
