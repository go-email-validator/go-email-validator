package view

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type SyntaxView struct {
	Username      string `json:"username"`
	Domain        string `json:"domain"`
	IsValidSyntax bool   `json:"is_valid_syntax"`
}

var NewSyntaxView = func(email email.EmailAddressInterface, result ev.ValidationResultInterface) ViewInterface {
	return ViewInterface(SyntaxView{
		email.Username(),
		email.Domain(),
		result.IsValid(),
	})
}
