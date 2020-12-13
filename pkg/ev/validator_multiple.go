package ev

import "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"

const MultipleValidatorName ValidatorName = "MultipleValidator"

type MultipleValidator struct {
	validators []ValidatorInterface
	AValidatorWithoutDeps
}

func (m MultipleValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var isValid = true
	var validator ValidatorInterface
	var vResult ValidationResultInterface

	for _, validator = range m.validators {
		vResult = validator.Validate(email)

		if !vResult.IsValid() {
			isValid = vResult.IsValid()
		}
	}

	return NewValidatorResult(isValid, nil, nil, "")
}

func NewMultipleValidator(validators ...ValidatorInterface) MultipleValidator {
	return MultipleValidator{validators: validators}
}
