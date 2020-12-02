package view

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type MXView struct {
	AcceptsMail bool     `json:"accepts_mail"`
	Records     []string `json:"records"`
}

var NewMXView = func(email email.EmailAddressInterface, result ev.ValidationResultInterface) ViewInterface {
	mxResult := result.(ev.MXValidationResultInterface)
	lenMX := len(mxResult.MX())
	records := make([]string, lenMX)

	for i, mx := range mxResult.MX() {
		records[i] = mx.Host
	}

	return ViewInterface(MXView{
		lenMX > 0,
		records,
	})
}
