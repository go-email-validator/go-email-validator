package converter

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"time"
)

type Name string

type Options interface {
	IsOptions()
	ExecutedTime() time.Duration
}

func NewOptions(executedTime time.Duration) Options {
	return options{
		ExecutedTimeValue: executedTime,
	}
}

type options struct {
	ExecutedTimeValue time.Duration
}

func (options) IsOptions() {}
func (o options) ExecutedTime() time.Duration {
	return o.ExecutedTimeValue
}

type Interface interface {
	Can(email evmail.Address, result ev.ValidationResult, opts Options) bool
	Convert(email evmail.Address, result ev.ValidationResult, opts Options) interface{}
}

type MapConverters map[ev.ValidatorName]Interface

func NewCompositeConverter(converters MapConverters) CompositeConverter {
	return CompositeConverter{converters}
}

type CompositeConverter struct {
	converters MapConverters
}

func (p CompositeConverter) converter(email evmail.Address, result ev.ValidationResult, opts Options) Interface {
	if converter, ok := p.converters[result.ValidatorName()]; ok && converter.Can(email, result, opts) {
		return converter
	}

	return nil
}

func (p CompositeConverter) Can(email evmail.Address, result ev.ValidationResult, opts Options) bool {
	return p.converter(email, result, opts) != nil
}

func (p CompositeConverter) Convert(email evmail.Address, result ev.ValidationResult, opts Options) interface{} {
	return p.converter(email, result, opts).Convert(email, result, opts)
}
