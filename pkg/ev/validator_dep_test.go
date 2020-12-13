package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testSleep struct {
	sleep time.Duration
	mockValidator
	deps    []ValidatorName
	results *[]ValidationResultInterface
}

func (t testSleep) GetDeps() []ValidatorName {
	return t.deps
}

func (t testSleep) Validate(_ ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	time.Sleep(t.sleep)

	var isValid = true
	for _, result := range *t.results {
		if !result.IsValid() {
			isValid = false
			break
		}
	}

	return NewValidatorResult(isValid && t.result, nil, nil, OtherValidator)
}

func TestDepValidator_Validate_Independent(t *testing.T) {
	email := getValidEmail()
	strings := emptyDeps

	depValidator := DepValidator{
		map[ValidatorName]ValidatorInterface{
			"test1": &testSleep{
				0,
				newMockValidator(true),
				strings,
				nil,
			},
			"test2": &testSleep{
				0,
				newMockValidator(true),
				strings,
				nil,
			},
			"test3": &testSleep{
				0,
				newMockValidator(false),
				strings,
				nil,
			},
		},
	}

	v := depValidator.Validate(email)
	assert.False(t, v.IsValid())
}

func TestDepValidator_Validate_Dependent(t *testing.T) {
	email := getValidEmail()
	strings := emptyDeps

	depValidator := DepValidator{
		map[ValidatorName]ValidatorInterface{
			"test1": &testSleep{
				100 * time.Millisecond,
				newMockValidator(true),
				strings,
				nil,
			},
			"test2": &testSleep{
				100 * time.Millisecond,
				newMockValidator(true),
				strings,
				nil,
			},
			"test3": &testSleep{
				100 * time.Millisecond,
				newMockValidator(true),
				[]ValidatorName{"test1", "test2"},
				nil,
			},
		},
	}

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}

func TestDepValidator_Validate_Full(t *testing.T) {
	email := ev_email.NewEmail("go.email.validator", "gmail.com")

	depValidator := DepValidator{
		map[ValidatorName]ValidatorInterface{
			RoleValidatorName:       NewRoleValidator(role.NewRBEASetRole()),
			DisposableValidatorName: NewDisposableValidator(disposable.MailCheckerDisposable{}),
			SyntaxValidatorName:     &SyntaxValidator{},
			MXValidatorName:         &MXValidator{},
			SMTPValidatorName: NewWarningsDecorator(
				ValidatorInterface(newSMTPValidator()),
				NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
					return func(err error) bool {
						return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
					}
				}),
			),
		},
	}

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
