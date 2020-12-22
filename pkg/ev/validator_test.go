package ev

import (
	email "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

var emptyEmail = email.NewEmailAddress("", "")
var validMockValidator Validator = mockValidator{result: true}
var inValidMockValidator Validator = mockValidator{result: false}

const validEmailString = "go.email.validator@gmail.com"

func newMockContains(value interface{}) mockContains {
	return mockContains{value}
}

type mockContains struct {
	value interface{}
}

func (m mockContains) Contains(value interface{}) bool {
	return m.value == value
}

func newMockError() error {
	return mockError{}
}

type mockError struct {
	utils.Err
}

func newMockValidator(result bool) mockValidator {
	return mockValidator{result: result}
}

type mockValidator struct {
	result bool
	AValidatorWithoutDeps
}

func (m mockValidator) Validate(_ email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	if !m.result {
		err = newMockError()
	}

	return NewValidatorResult(m.result, utils.Errs(err), nil, OtherValidator)
}

func TestMockValidator(t *testing.T) {
	cases := []struct {
		validator mockValidator
		expected  ValidationResult
	}{
		{
			validator: newMockValidator(true),
			expected:  NewValidatorResult(true, nil, nil, OtherValidator),
		},
		{
			validator: newMockValidator(false),
			expected:  NewValidatorResult(false, utils.Errs(newMockError()), nil, OtherValidator),
		},
	}

	var emptyEmail email.EmailAddress
	for _, c := range cases {
		actual := c.validator.Validate(emptyEmail)
		assert.Equal(t, c.expected, actual)
	}
}

func TestAValidatorWithoutDeps(t *testing.T) {
	validator := AValidatorWithoutDeps{}

	assert.Equal(t, emptyDeps, validator.GetDeps())
}
