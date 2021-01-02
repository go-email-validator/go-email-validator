package ev

import (
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

// Create warning decorator to skip some errors
func NewWarningsDecorator(validator Validator, isWarning IsWarning) Validator {
	return warningsDecorator{validator, isWarning}
}

// Detect error as warning
type IsWarning func(err error) bool

type WarningSet sets.Set

// Create function for detection of warnings
func NewIsWarning(warningMap WarningSet, isWarning func(warningMap WarningSet) IsWarning) IsWarning {
	return isWarning(warningMap)
}

type warningsDecorator struct {
	validator Validator
	isWarning IsWarning
}

func (w warningsDecorator) GetDeps() []ValidatorName {
	return w.validator.GetDeps()
}

func (w warningsDecorator) Validate(email evmail.Address, results ...ValidationResult) ValidationResult {
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
