package ev

import (
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockValidator struct {
	result bool
}

func (m mockValidator) Validate(_ email.EmailAddressInterface) ValidationResultInterface {
	return NewValidatorResult(m.result, nil, nil)
}

var validEmail email.EmailAddressInterface = email.EmailAddress{}
var validMockValidator ValidatorInterface = mockValidator{true}
var inValidMockValidator ValidatorInterface = mockValidator{false}

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
