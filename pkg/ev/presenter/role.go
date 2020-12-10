package presenter

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type RolePresenter struct {
	IsRoleAccount bool `json:"is_role_account"`
}

type RoleProcessor struct{}

func (s RoleProcessor) CanProcess(_ email.EmailAddressInterface, result ev.ValidationResultInterface) bool {
	return result.ValidatorName() == ev.RoleValidatorName
}

func (s RoleProcessor) Process(email email.EmailAddressInterface, result ev.ValidationResultInterface) interface{} {
	return RolePresenter{result.IsValid()}
}
