package presentation

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"time"
)

type ValidationPresenter interface {
	Validate(email string, name converter.Name, opts map[ev.ValidatorName]interface{}) (interface{}, error)
}

func NewValidationPresenter(presenters map[converter.Name]Interface) ValidationPresenter {
	return multiplePresenter{presenters: presenters}
}

type multiplePresenter struct {
	presenters map[converter.Name]Interface
}

func (p multiplePresenter) Validate(email string, name converter.Name, opts map[ev.ValidatorName]interface{}) (interface{}, error) {
	_presenter, ok := p.presenters[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("converter with name \"%s\" does not exist", name))
	}

	return _presenter.Validate(email, opts)
}

type GetEmail func(email string) evmail.Address

type Interface interface {
	Validate(email string, opts map[ev.ValidatorName]interface{}) (interface{}, error)
	ValidateFromAddress(email evmail.Address, opts map[ev.ValidatorName]interface{}) (interface{}, error)
}

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
