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

// NewInput create Input from evmail.Address and KVOption list
func NewInput(email evmail.Address, kvOptions ...KVOption) Input {
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

// NewKVOption instantiates KVOption
func NewKVOption(name ValidatorName, option interface{}) KVOption {
	return KVOption{
		Name:   name,
		Option: option,
	}
}

// KVOption needs to form options in Input
type KVOption struct {
	Name   ValidatorName
	Option interface{}
}
