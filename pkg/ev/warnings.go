package ev

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

func NewWarningsDecorator(validator Validator, isWarning IsWarning) Validator {
	return WarningsDecorator{validator, isWarning}
}

type IsWarning func(err error) bool

type WarningSet sets.Set

func NewIsWarning(warningMap WarningSet, isWarning func(warningMap WarningSet) IsWarning) IsWarning {
	return isWarning(warningMap)
}

type WarningsDecorator struct {
	validator Validator
	isWarning IsWarning
}

func (w WarningsDecorator) GetDeps() []ValidatorName {
	return w.validator.GetDeps()
}

func (w WarningsDecorator) Validate(email evmail.Address, results ...ValidationResult) ValidationResult {
	result := w.validator.Validate(email, results...)
	changeableResult, ok := result.(ChangeableValidationResult)
	if !ok {
		return result
	}

	var errors, warnings []error
	for _, err := range result.Errors() {
		if w.isWarning(err) {
			warnings = append(warnings, err)
		} else {
			errors = append(errors, err)
		}
	}

	changeableResult.SetErrors(errors)
	changeableResult.SetWarnings(warnings)

	return changeableResult.(ValidationResult)
}
