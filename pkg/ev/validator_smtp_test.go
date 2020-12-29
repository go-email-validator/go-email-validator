package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/stretchr/testify/assert"
	"testing"
)

// test monicaramirezrestrepo@hotmail.com.
func newSMTPValidator() Validator {
	return NewSMTPValidator(evsmtp.NewChecker(evsmtp.CheckerDTO{
		DialFunc:  evsmtp.Dial,
		SendMail:  evsmtp.NewSendMail(),
		FromEmail: evmail.FromString(evsmtp.DefaultEmail),
	}))
}

func getSmtpValidator_Validate() Validator {
	return NewDepValidator(
		map[ValidatorName]Validator{
			SyntaxValidatorName: NewSyntaxValidator(),
			MXValidatorName:     DefaultNewMXValidator(),
			SMTPValidatorName: NewWarningsDecorator(
				Validator(newSMTPValidator()),
				NewIsWarning(hashset.New(evsmtp.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
					return func(err error) bool {
						return warningMap.Contains(err.(evsmtp.Error).Stage())
					}
				}),
			),
		},
	)
}

func BenchmarkSMTPValidator_Validate(b *testing.B) {
	email := evmail.FromString(validEmailString)
	validator := getSmtpValidator_Validate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(email)
	}
}

func TestSMTPValidator_Validate_WithoutMock(t *testing.T) {
	email := evmail.FromString(validEmailString)
	validator := getSmtpValidator_Validate()

	v := validator.Validate(email)
	assert.True(t, v.IsValid())
}
