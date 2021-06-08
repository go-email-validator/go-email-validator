package converter

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"time"
)

// Name of converter
type Name string

// Options changes the process of converting
type Options interface {
	IsOptions()
	ExecutedTime() time.Duration
}

// NewOptions creates Options
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

// Interface converts ev.ValidationResult in some presenter
type Interface interface {
	// Can defines the possibility applying of a converter
	Can(email evmail.Address, result ev.ValidationResult, opts Options) bool
	// Convert converts ev.ValidationResult in some presenter
	Convert(email evmail.Address, result ev.ValidationResult, opts Options) interface{}
}

// MapConverters is a map of converters
type MapConverters map[ev.ValidatorName]Interface

// NewCompositeConverter creates CompositeConverter
func NewCompositeConverter(converters MapConverters) CompositeConverter {
	return CompositeConverter{converters}
}

// CompositeConverter converts ev.ValidationResult depends of ev.ValidationResult.ValidatorName()
type CompositeConverter struct {
	converters MapConverters
}

func (p CompositeConverter) converter(email evmail.Address, result ev.ValidationResult, opts Options) Interface {
	if converter, ok := p.converters[result.ValidatorName()]; ok && converter.Can(email, result, opts) {
		return converter
	}

	return nil
}

// Can result ev.ValidationResult be converted
func (p CompositeConverter) Can(email evmail.Address, result ev.ValidationResult, opts Options) bool {
	return p.converter(email, result, opts) != nil
}

// Convert ev.ValidationResult depends of ev.ValidationResult.ValidatorName()
func (p CompositeConverter) Convert(email evmail.Address, result ev.ValidationResult, opts Options) interface{} {
	return p.converter(email, result, opts).Convert(email, result, opts)
}
