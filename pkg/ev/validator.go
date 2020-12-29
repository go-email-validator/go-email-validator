package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

type ValidatorName string

const OtherValidator ValidatorName = "other"

func (v ValidatorName) String() string {
	return string(v)
}

type Validator interface {
	GetDeps() []ValidatorName
	Validate(email evmail.Address, results ...ValidationResult) ValidationResult
}

type ChangeableValidationResult interface {
	SetErrors([]error)
	SetWarnings([]error)
}

type ValidationResult interface {
	IsValid() bool
	Errors() []error
	HasErrors() bool
	Warnings() []error
	HasWarnings() bool
	ValidatorName() ValidatorName
}

// Abstract class for expected of validation
type AValidationResult struct {
	isValid  bool
	errors   []error
	warnings []error
	name     ValidatorName
}

func (a AValidationResult) IsValid() bool {
	return a.isValid
}

func (a *AValidationResult) SetErrors(errors []error) {
	a.errors = errors
	a.isValid = !a.HasErrors()
}

func (a *AValidationResult) Errors() []error {
	return a.errors
}

func (a *AValidationResult) HasErrors() bool {
	return utils.RangeLen(a.Errors()) > 0
}
func (a *AValidationResult) SetWarnings(warnings []error) {
	a.warnings = warnings
}
func (a *AValidationResult) Warnings() []error {
	return a.warnings
}

func (a *AValidationResult) HasWarnings() bool {
	return utils.RangeLen(a.Warnings()) > 0
}

func (a AValidationResult) ValidatorName() ValidatorName {
	return a.name
}

type validationResult = AValidationResult

func NewValidValidatorResult(name ValidatorName) ValidationResult {
	return NewValidatorResult(true, nil, nil, name)
}

func NewValidatorResult(isValid bool, errors []error, warnings []error, name ValidatorName) ValidationResult {
	if name == "" {
		name = OtherValidator
	}

	return &validationResult{isValid, errors, warnings, name}
}

var emptyDeps = make([]ValidatorName, 0)

type AValidatorWithoutDeps struct{}

func (a AValidatorWithoutDeps) GetDeps() []ValidatorName {
	return emptyDeps
}
