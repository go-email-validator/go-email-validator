package ev

import (
	ev_email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestADepValidator_SetResults(t *testing.T) {
	v1 := ADepValidator{nil}
	v2 := v1.SetResults(ValidationResult{})

	v22 := v2.(ADepValidator)

	assert.Equal(t, v1, v22)
}

type type1 struct{}
type type2 struct{}

func TestGetDeps(t *testing.T) {
	r := GetDepNames((*type1)(nil), (*type2)(nil))

	assert.Equal(t, r, []string{"type1", "type2"})
}

type testSleep struct {
	sleep time.Duration
	mockValidator
	deps    []string
	results *[]ValidationResultInterface
}

func (t testSleep) GetDeps() []string {
	return t.deps
}
func (t *testSleep) SetResults(results ...ValidationResultInterface) ValidatorInterface {
	t.results = &results

	return t
}

func (t testSleep) Validate(_ ev_email.EmailAddressInterface) ValidationResultInterface {
	time.Sleep(t.sleep)

	var isValid = true
	for _, result := range *t.results {
		if !result.IsValid() {
			isValid = false
			break
		}
	}

	return NewValidatorResult(isValid && t.result, nil, nil)
}

func getEmptyDeps() []string {
	return []string{}
}

func TestDepValidator_Validate_Independent(t *testing.T) {
	email := getValidEmail()
	strings := getEmptyDeps()

	depValidator := DepValidator{
		map[string]ValidatorInterface{
			"test1": &testSleep{
				0,
				mockValidator{true},
				strings,
				nil,
			},
			"test2": &testSleep{
				0,
				mockValidator{true},
				strings,
				nil,
			},
			"test3": &testSleep{
				0,
				mockValidator{false},
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
	strings := getEmptyDeps()

	depValidator := DepValidator{
		map[string]ValidatorInterface{
			"test1": &testSleep{
				100 * time.Millisecond,
				mockValidator{true},
				strings,
				nil,
			},
			"test2": &testSleep{
				100 * time.Millisecond,
				mockValidator{true},
				strings,
				nil,
			},
			"test3": &testSleep{
				100 * time.Millisecond,
				mockValidator{true},
				[]string{"test1", "test2"},
				nil,
			},
		},
	}

	v := depValidator.Validate(email)
	assert.True(t, v.IsValid())
}
