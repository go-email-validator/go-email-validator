package presenter

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type DisposablePresenter struct {
	IsDisposable bool `json:"is_disposable"`
}

type DisposableProcessor struct{}

func (s DisposableProcessor) CanProcess(_ email.EmailAddressInterface, result ev.ValidationResultInterface) bool {
	return result.ValidatorName() == ev.DisposableValidatorName
}

func (s DisposableProcessor) Process(email email.EmailAddressInterface, result ev.ValidationResultInterface) interface{} {
	return DisposablePresenter{result.IsValid()}
}
