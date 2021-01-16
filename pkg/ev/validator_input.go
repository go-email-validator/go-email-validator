package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

// ValidatorName is type to represent validator name
type ValidatorName string

func (v ValidatorName) String() string {
	return string(v)
}

type Interface interface {
	Email() evmail.Address
	Option(name ValidatorName) interface{}
}

// NewInput create Interface from evmail.Address and KVOption list
func NewInput(email evmail.Address, kvOptions ...KVOption) Interface {
	var options = make(map[ValidatorName]interface{})

	for _, kvOption := range kvOptions {
		options[kvOption.Name] = kvOption.Option
	}

	return NewInputFromMap(email, options)
}

// NewInputFromMap create Interface from evmail.Address and options
func NewInputFromMap(email evmail.Address, options map[ValidatorName]interface{}) Interface {
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

func NewKVOption(name ValidatorName, option interface{}) KVOption {
	return KVOption{
		Name:   name,
		Option: option,
	}
}

type KVOption struct {
	Name   ValidatorName
	Option interface{}
}
