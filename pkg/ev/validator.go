package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
)

type ValidatorInterface interface {
	GetDeps() []string
	Validate(email ev_email.EmailAddressInterface, results ...ValidationResultInterface) ValidationResultInterface
}

type ChangeableValidationResultInterface interface {
	SetErrors([]error)
	SetWarnings([]error)
}

type ValidationResultInterface interface {
	IsValid() bool
	Errors() []error
	HasErrors() bool
	Warnings() []error
	HasWarnings() bool
}

// Abstract class for result of validation
type AValidationResult struct {
	isValid  bool
	errors   []error
	warnings []error
}

func (a AValidationResult) IsValid() bool {
	return a.isValid
}

func (a *AValidationResult) SetErrors(errors []error) {
	a.isValid = len(errors) == 0
	a.errors = errors
}

func (a AValidationResult) Errors() []error {
	return a.errors
}

func (a AValidationResult) HasErrors() bool {
	return utils.RangeLen(a.Errors()) > 0
}
func (a *AValidationResult) SetWarnings(warnings []error) {
	a.warnings = warnings
}
func (a AValidationResult) Warnings() []error {
	return a.warnings
}

func (a AValidationResult) HasWarnings() bool {
	return utils.RangeLen(a.Warnings()) > 0
}

type ValidationResult = AValidationResult

func NewValidatorResult(isValid bool, errors []error, warnings []error) ValidationResultInterface {
	return &ValidationResult{isValid, errors, warnings}
}

var emptyStrings = make([]string, 0)

type AValidatorWithoutDeps struct {
}

func (A AValidatorWithoutDeps) GetDeps() []string {
	return emptyStrings
}
