package ev

type ValidatorInterface interface {
	Validate(email EmailAddressInterface) ValidationResultInterface
}

type ValidationResultInterface interface {
	IsValid() bool
	Errors() interface{}
	HasErrors() bool
	Warnings() interface{}
	HasWarnings() bool
}

// Abstract class for result of validation
type AValidationResult struct {
	isValid  bool
	errors   interface{}
	warnings interface{}
}

func (a AValidationResult) IsValid() bool {
	return a.isValid
}

func (a AValidationResult) Errors() interface{} {
	return a.errors
}

func (a AValidationResult) HasErrors() bool {
	return RangeLen(a.Errors()) > 0
}

func (a AValidationResult) Warnings() interface{} {
	return a.warnings
}

func (a AValidationResult) HasWarnings() bool {
	return RangeLen(a.Warnings()) > 0
}

type ValidationResult = AValidationResult

func NewValidatorResult( /*t reflect.Type,*/ isValid bool, errors interface{}, warnings interface{}) ValidationResultInterface {
	/*var validatorResult ValidationResultInterface
	if t == nil {
		validatorResult = new(ValidationResult)
	} else {
		validatorResult = reflect.New(t).Interface().(ValidationResultInterface)
	}

	return validatorResult*/

	return ValidationResult{isValid, errors, warnings}
}
