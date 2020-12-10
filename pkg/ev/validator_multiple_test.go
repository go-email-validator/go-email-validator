package ev

import (
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockValidator(result bool) mockValidator {
	return mockValidator{result: result}
}

type mockValidator struct {
	result bool
	AValidatorWithoutDeps
}

func (m mockValidator) Validate(_ email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	return NewValidatorResult(m.result, nil, nil, OtherValidator)
}

var validEmail email.EmailAddressInterface = email.EmailAddress{}
var validMockValidator ValidatorInterface = mockValidator{result: true}
var inValidMockValidator ValidatorInterface = mockValidator{result: false}

func TestMultipleValidatorInValid(t *testing.T) {
	var v ValidatorInterface = NewMultipleValidator(validMockValidator, inValidMockValidator)

	assert.False(
		t,
		v.Validate(validEmail).IsValid(),
		"should be invalid",
	)
}

func TestMultipleValidatorValid(t *testing.T) {
	v := NewMultipleValidator(validMockValidator, validMockValidator)

	assert.True(
		t,
		v.Validate(validEmail).IsValid(),
		"should be valid",
	)
}
