package ev

import "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"

const MultipleValidatorName = "MultipleValidatorInterface"

type MultipleValidator struct {
	validators []ValidatorInterface
}

func (m MultipleValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	var isValid = true
	var validator ValidatorInterface
	var vResult ValidationResultInterface

	for _, validator = range m.validators {
		vResult = validator.Validate(email)

		if !vResult.IsValid() {
			isValid = vResult.IsValid()
		}
	}

	return NewValidatorResult(isValid, nil, nil)
}

func NewMultipleValidator(validators ...ValidatorInterface) MultipleValidator {
	return MultipleValidator{validators}
}
