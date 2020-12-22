package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"testing"
)

func BenchmarkSMTPValidator_Validate_MX(b *testing.B) {
	email := ev_email.EmailFromString(validEmailString)

	depValidator := NewDepValidator(
		map[ValidatorName]Validator{
			SyntaxValidatorName: NewMXValidator(),
			MXValidatorName:     NewSyntaxValidator(),
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}
