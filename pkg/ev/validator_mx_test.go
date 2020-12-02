package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"testing"
)

func BenchmarkSMTPValidator_Validate_MX(b *testing.B) {
	email := ev_email.NewEmail("go.email.validator", "gmail.com")

	depValidator := DepValidator{
		map[string]ValidatorInterface{
			SyntaxValidatorName: &SyntaxValidator{},
			MXValidatorName:     &MXValidator{},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}
