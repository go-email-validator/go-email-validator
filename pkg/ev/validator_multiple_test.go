package ev

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockValidator struct {
	result bool
}

func (m mockValidator) Validate(email EmailAddressInterface) ValidationResultInterface {
	return NewValidatorResult(m.result, nil, nil)
}

var validEmail EmailAddressInterface = EmailAddress{}
var validMockValidator ValidatorInterface = mockValidator{true}
var inValidMockValidator ValidatorInterface = mockValidator{false}

func TestMultipleValidatorInValid(t *testing.T) {
	v := NewMultipleValidator(validMockValidator, inValidMockValidator)

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
