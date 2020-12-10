package presenter

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type MXPresenter struct {
	AcceptsMail bool     `json:"accepts_mail"`
	Records     []string `json:"records"`
}

type MXProcessor struct{}

func (s MXProcessor) CanProcess(_ email.EmailAddressInterface, result ev.ValidationResultInterface) bool {
	return result.ValidatorName() == ev.MXValidatorName
}

func (s MXProcessor) Process(email email.EmailAddressInterface, result ev.ValidationResultInterface) interface{} {
	mxResult := result.(ev.MXValidationResultInterface)
	lenMX := len(mxResult.MX())
	records := make([]string, lenMX)

	for i, mx := range mxResult.MX() {
		records[i] = mx.Host
	}

	return MXPresenter{
		lenMX > 0,
		records,
	}
}
