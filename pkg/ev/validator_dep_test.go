package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
	"github.com/stretchr/testify/assert"
	"net/smtp"
	"testing"
	"time"
)

type testSleep struct {
	sleep time.Duration
	mockValidator
	deps []ValidatorName
}

func (t testSleep) GetDeps() []ValidatorName {
	return t.deps
}

func (t testSleep) Validate(_ ev_email.EmailAddress, results ...ValidationResult) ValidationResult {
	time.Sleep(t.sleep)

	var isValid = true
	for _, result := range results {
		if !result.IsValid() {
			isValid = false
			break
		}
	}

	return NewValidatorResult(isValid && t.result, nil, nil, OtherValidator)
}

func TestDepValidator_Validate_Independent(t *testing.T) {
	email := GetValidTestEmail()
	strings := emptyDeps

	depValidator := NewDepValidator(
		map[ValidatorName]Validator{
			"test1": &testSleep{
				0,
				newMockValidator(true),
				strings,
			},
			"test2": &testSleep{
				0,
				newMockValidator(true),
				strings,
			},
			"test3": &testSleep{
				0,
				newMockValidator(false),
				strings,
			},
		},
	)

	v := depValidator.Validate(email)
	assert.False(t, v.IsValid())
}

func TestDepValidator_Validate_Dependent(t *testing.T) {
	email := GetValidTestEmail()
	strings := emptyDeps

	depValidator := NewDepValidator(map[ValidatorName]Validator{
		"test1": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			strings,
		},
		"test2": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			strings,
		},
		"test3": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			[]ValidatorName{"test1", "test2"},
		},
	},
	)

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}

func TestDepValidator_Validate_Full(t *testing.T) {
	email := ev_email.EmailFromString(validEmailString)

	depValidator := NewDepValidator(map[ValidatorName]Validator{
		//FreeValidatorName:     FreeDefaultValidator(),
		RoleValidatorName:       NewRoleValidator(role.NewRBEASetRole()),
		DisposableValidatorName: NewDisposableValidator(contains.NewFunc(disposable.MailChecker)),
		SyntaxValidatorName:     NewSyntaxValidator(),
		MXValidatorName:         NewMXValidator(),
		SMTPValidatorName: NewWarningsDecorator(
			NewSMTPValidator(smtp_checker.Checker{
				DialFunc:  smtp.Dial,
				SendMail:  smtp_checker.NewSendMail(),
				FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
			}),
			NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
				return func(err error) bool {
					return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
				}
			}),
		),
	},
	)

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
