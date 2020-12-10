package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
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
