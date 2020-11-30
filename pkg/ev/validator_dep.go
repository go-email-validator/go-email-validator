package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"reflect"
	"sync"
)

type DepValidatorInterface interface {
	ValidatorInterface
	// Use GetDepNames to get slice
	GetDeps() []string
	SetResults(results ...ValidationResultInterface) ValidatorInterface
}

type ADepValidator struct {
	results *[]ValidationResultInterface
}

func (a *ADepValidator) GetDeps() []string {
	panic("implement me")
}

func (a *ADepValidator) SetResults(results ...ValidationResultInterface) ValidatorInterface {
	a.results = &results
	return a
}

func (a ADepValidator) Validate(_ ev_email.EmailAddressInterface) ValidationResultInterface {
	panic("implement me")
}

// dep - (*struct)(nil)
func GetDepNames(deps ...interface{}) []string {
	var result = make([]string, len(deps))

	for i, dep := range deps {
		result[i] = GetDepName(dep)
	}

	return result
}

func GetDepName(dep interface{}) string {
	return reflect.TypeOf(dep).Elem().Name()
}

func NewDepValidator(deps map[string]ValidatorInterface) DepValidator {
	return DepValidator{deps}
}

type DepValidator struct {
	deps map[string]ValidatorInterface
}

func (d *DepValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	var waiters, waitersMutex = make(map[string][]*sync.WaitGroup), sync.RWMutex{}
	var validationResults, validationResultsMutex = make(map[string]ValidationResultInterface), sync.RWMutex{}
	var isValid = true
	var starter, finisher = sync.WaitGroup{}, sync.WaitGroup{}
	starter.Add(1)
	finisher.Add(len(d.deps))

	for key, validator := range d.deps {
		var depWaiter *sync.WaitGroup
		var depWaiters []*sync.WaitGroup
		var deps []string

		v, ok := validator.(DepValidatorInterface)
		if ok {
			deps = v.GetDeps()

			depWaiter = &sync.WaitGroup{}
			depWaiter.Add(len(deps))

			for _, dep := range deps {
				if depWaiters, ok = waiters[dep]; !ok {
					depWaiters = make([]*sync.WaitGroup, 0)
				}

				waiters[dep] = append(depWaiters, depWaiter)
			}
		}

		go func(key string, validator ValidatorInterface, depWaiter *sync.WaitGroup) {
			// add recover
			starter.Wait()
			if depWaiter != nil {
				depWaiter.Wait()

				var results = make([]ValidationResultInterface, len(deps))
				validationResultsMutex.RLock()
				for i, dep := range deps {
					results[i] = validationResults[dep]
				}
				validationResultsMutex.RUnlock()
				validator.(DepValidatorInterface).SetResults(results...)
			}

			var result = validator.Validate(email)
			validationResultsMutex.Lock()
			validationResults[key] = result
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

	return NewValidatorResult(isValid, nil, nil)
}
