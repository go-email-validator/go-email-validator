package ev

import (
	email "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
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
