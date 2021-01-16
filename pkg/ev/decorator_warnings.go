package ev

import (
	"github.com/emirpasic/gods/sets"
)

// NewWarningsDecorator creates warning decorator to skip some errors
func NewWarningsDecorator(validator Validator, isWarning IsWarning) Validator {
	return warningsDecorator{validator, isWarning}
}

// IsWarning is type to detect error as warning
type IsWarning func(err error) bool

// WarningSet is alias for sets.Set
type WarningSet sets.Set

// NewIsWarning creates function for detection of warnings
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

func (w warningsDecorator) Validate(input Interface, results ...ValidationResult) ValidationResult {
	result := w.validator.Validate(input, results...)
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
