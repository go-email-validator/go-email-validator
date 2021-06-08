package presentation

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"time"
)

// GetEmail converts string to evmail.Address
type GetEmail func(email string) evmail.Address

// Interface is a decorator to represent result in a form of some service
type Interface interface {
	Validate(email string, opts map[ev.ValidatorName]interface{}) (interface{}, error)
	ValidateFromAddress(email evmail.Address, opts map[ev.ValidatorName]interface{}) (interface{}, error)
}

// NewPresenter creates presenter decorator
func NewPresenter(getEmail GetEmail, validator ev.Validator, converter converter.Interface) Interface {
	return presenter{
		getEmail:  getEmail,
		validator: validator,
		converter: converter,
	}
}

type presenter struct {
	getEmail  func(email string) evmail.Address
	validator ev.Validator
	converter converter.Interface
}

func (p presenter) Validate(email string, opts map[ev.ValidatorName]interface{}) (interface{}, error) {
	address := p.getEmail(email)

	return p.ValidateFromAddress(address, opts)
}

func (p presenter) ValidateFromAddress(address evmail.Address, opts map[ev.ValidatorName]interface{}) (interface{}, error) {
	start := time.Now()
	validationResult := p.validator.Validate(ev.NewInputFromMap(address, opts))
	elapsed := time.Since(start)

	return p.converter.Convert(address, validationResult, converter.NewOptions(elapsed)), nil
}
