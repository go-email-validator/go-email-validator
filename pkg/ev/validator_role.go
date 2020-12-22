package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const RoleValidatorName ValidatorName = "RoleValidator"

type RoleError struct {
	utils.Err
}

func NewRoleValidator(r contains.Interface) Validator {
	return roleValidator{r: r}
}

type roleValidator struct {
	r contains.Interface
	AValidatorWithoutDeps
}

func (r roleValidator) Validate(email ev_email.EmailAddress, _ ...ValidationResult) ValidationResult {
	var err error
	var hasRole = r.r.Contains(email.Username())
	if hasRole {
		err = RoleError{}
	}

	return NewValidatorResult(!hasRole, utils.Errs(err), nil, RoleValidatorName)
}
