package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	mockevmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"github.com/stretchr/testify/require"
	"testing"
)

// test monicaramirezrestrepo@hotmail.com.
func newSMTPValidator() Validator {
	return NewSMTPValidator(evsmtp.NewChecker(evsmtp.CheckerDTO{
		DialFunc: evsmtp.DirectDial,
		SendMail: evsmtp.NewSendMail(nil),
		Options: evsmtp.NewOptions(evsmtp.OptionsDTO{
			EmailFrom: evmail.FromString(evsmtp.DefaultEmail),
		}),
	}))
}

func getSMTPValidatorValidate() Validator {
	return NewDepValidator(
		map[ValidatorName]Validator{
			SyntaxValidatorName: NewSyntaxValidator(),
			MXValidatorName:     DefaultNewMXValidator(),
			SMTPValidatorName: NewWarningsDecorator(
				newSMTPValidator(),
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
	email := evmail.FromString(mockevmail.ValidEmailString)
	validator := getSMTPValidatorValidate()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(NewInput(email))
	}
}

func TestSMTPValidator_Validate_WithoutMock(t *testing.T) {
	email := evmail.FromString(mockevmail.ValidEmailString)
	validator := getSMTPValidatorValidate()

	v := validator.Validate(NewInput(email))
	require.True(t, v.IsValid())
}
