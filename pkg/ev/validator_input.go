package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

// ValidatorName is type to represent validator name
type ValidatorName string

func (v ValidatorName) String() string {
	return string(v)
}

// Input consists of input data for Validator.Validate
type Input interface {
	Email() evmail.Address
	Option(name ValidatorName) interface{}
}

// NewInput create Input from evmail.Address and kvOption list
func NewInput(email evmail.Address, kvOptions ...kvOption) Input {
	var options = make(map[ValidatorName]interface{})

	for _, kvOption := range kvOptions {
		options[kvOption.Name] = kvOption.Option
	}

	return NewInputFromMap(email, options)
}

// NewInputFromMap create Input from evmail.Address and options
func NewInputFromMap(email evmail.Address, options map[ValidatorName]interface{}) Input {
	return &input{
		email:   email,
		options: options,
	}
}

type input struct {
	email   evmail.Address
	options map[ValidatorName]interface{}
}

func (i *input) Email() evmail.Address {
	return i.email
}

func (i *input) Option(name ValidatorName) interface{} {
	return i.options[name]
}

// KVOption instantiates kvOption
func KVOption(name ValidatorName, option interface{}) kvOption {
	return kvOption{
		Name:   name,
		Option: option,
	}
}

// kvOption needs to form options in Input
type kvOption struct {
	Name   ValidatorName
	Option interface{}
}
