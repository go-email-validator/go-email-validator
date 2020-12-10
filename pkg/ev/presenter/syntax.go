package presenter

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type SyntaxPresenter struct {
	Username      string `json:"username"`
	Domain        string `json:"domain"`
	IsValidSyntax bool   `json:"is_valid_syntax"`
}

type SyntaxProcessor struct{}

func (s SyntaxProcessor) CanProcess(_ email.EmailAddressInterface, result ev.ValidationResultInterface) bool {
	return result.ValidatorName() == ev.SyntaxValidatorName
}

func (s SyntaxProcessor) Process(email email.EmailAddressInterface, result ev.ValidationResultInterface) interface{} {
	return SyntaxPresenter{
		email.Username(),
		email.Domain(),
		result.IsValid(),
	}
}
