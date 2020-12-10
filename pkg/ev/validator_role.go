package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/role"
)

const RoleValidatorName ValidatorName = "RoleValidator"

func NewRoleValidator(r role.Interface) ValidatorInterface {
	return RoleValidator{r: r}
}

type RoleValidator struct {
	r role.Interface
	AValidatorWithoutDeps
}

func (r RoleValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	return NewValidatorResult(r.r.HasRole(email), nil, nil, RoleValidatorName)
}
