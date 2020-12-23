package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"sync"
)

const DepValidatorName ValidatorName = "depValidator"

type ValidatorMap map[ValidatorName]Validator

type DepsError struct {
	utils.Err
}

func NewDepValidator(deps ValidatorMap) Validator {
	return depValidator{deps: deps}
}

type depValidator struct {
	AValidatorWithoutDeps
	deps ValidatorMap
}

func (d depValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var waiters, waitersMutex = make(map[ValidatorName][]*sync.WaitGroup), sync.RWMutex{}
	var validationResultsByName, validationResultsMutex = make(map[ValidatorName]ValidationResult), sync.RWMutex{}
	var isValid = true
	var starter, finisher = sync.WaitGroup{}, sync.WaitGroup{}
	starter.Add(1)
	finisher.Add(len(d.deps))

	for key, validator := range d.deps {
		var depWaiter *sync.WaitGroup
		var depWaiters []*sync.WaitGroup
		var deps []ValidatorName
		var ok bool

		deps = validator.GetDeps()
		if len(deps) > 0 {
			depWaiter = &sync.WaitGroup{}
			depWaiter.Add(len(deps))

			for _, dep := range deps {
				if depWaiters, ok = waiters[dep]; !ok {
					depWaiters = make([]*sync.WaitGroup, 0)
				}
				waiters[dep] = append(depWaiters, depWaiter)
			}
		}

		go func(key ValidatorName, validator Validator, depWaiter *sync.WaitGroup) {
			var results []ValidationResult

			// TODO add recover
			starter.Wait()
			if depWaiter != nil {
				depWaiter.Wait()

				results = make([]ValidationResult, len(deps))
				validationResultsMutex.RLock()
				for i, dep := range deps {
					results[i] = validationResultsByName[dep]
				}
				validationResultsMutex.RUnlock()
			}

			var result = validator.Validate(email, results...)
			validationResultsMutex.Lock()
			validationResultsByName[key] = result
			isValid = isValid && result.IsValid()
			validationResultsMutex.Unlock()

			waitersMutex.RLock()
			if depWaiters, ok = waiters[key]; ok {
				for _, depWaiter := range depWaiters {
					depWaiter.Done()
				}
			}
			waitersMutex.RUnlock()
			finisher.Done()
		}(key, validator, depWaiter)
	}
	starter.Done()
	finisher.Wait()

	return NewDepValidatorResult(isValid, validationResultsByName)
}

type DepResult map[ValidatorName]ValidationResult

type DepValidationResult interface {
	ValidationResult
	GetResults() DepResult
}

func NewDepValidatorResult(isValid bool, results DepResult) ValidationResult {
	return depValidationResult{
		isValid,
		results,
	}
}

type depValidationResult struct {
	isValid bool
	results DepResult
}

func (d depValidationResult) GetResults() DepResult {
	return d.results
}

func (d depValidationResult) IsValid() bool {
	return d.isValid
}

func (d depValidationResult) Errors() []error {
	var errors = make([]error, 0)

	for _, result := range d.GetResults() {
		for _, err := range result.Errors() {
			errors = append(errors, err)
		}
	}

	return errors
}

func (d depValidationResult) HasErrors() bool {
	for _, result := range d.GetResults() {
		if result.HasErrors() {
			return true
		}
	}

	return false
}

func (d depValidationResult) Warnings() []error {
	var warnings = make([]error, 0)

	for _, result := range d.GetResults() {
		for _, warning := range result.Warnings() {
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

func (d depValidationResult) HasWarnings() bool {
	for _, result := range d.GetResults() {
		if result.HasWarnings() {
			return true
		}
	}

	return false
}

func (d depValidationResult) ValidatorName() ValidatorName {
	return DepValidatorName
}
