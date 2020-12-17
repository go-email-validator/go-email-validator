package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
)

const RoleValidatorName ValidatorName = "RoleValidator"

type RoleError struct {
	error
}

func NewRoleValidator(r role.Interface) ValidatorInterface {
	return RoleValidator{r: r}
}

type RoleValidator struct {
	r role.Interface
	AValidatorWithoutDeps
}

func (r RoleValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var errs = make([]error, 0)
	var hasRole = r.r.HasRole(email)
	if hasRole {
		errs = append(errs, RoleError{})
	}

	return NewValidatorResult(!hasRole, errs, nil, RoleValidatorName)
}
