package checkifemailexist

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

type disposablePresentation struct {
	IsDisposable bool `json:"is_disposable"`
}

type disposableConverter struct{}

func (disposableConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.DisposableValidatorName
}

func (disposableConverter) Convert(_ evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	return disposablePresentation{!result.IsValid()}
}
