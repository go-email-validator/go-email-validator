package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/role"
)

const RoleValidatorName = "RoleValidatorInterface"

type RoleValidatorInterface interface {
	ValidatorInterface
}

func NewRoleValidator(r role.Interface) ValidatorInterface {
	return RoleValidator{r}
}

type RoleValidator struct {
	r role.Interface
}

func (r RoleValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	return NewValidatorResult(r.r.HasRole(email), nil, nil)
}
