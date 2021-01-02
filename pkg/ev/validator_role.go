package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

const RoleValidatorName ValidatorName = "RoleValidator"

type RoleError struct {
	utils.Err
}

func NewRoleValidator(r contains.InSet) Validator {
	return roleValidator{r: r}
}

type roleValidator struct {
	AValidatorWithoutDeps
	r contains.InSet
}

func (r roleValidator) Validate(email evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	var hasRole = r.r.Contains(email.Username())
	if hasRole {
		err = RoleError{}
	}

	return NewResult(!hasRole, utils.Errs(err), nil, RoleValidatorName)
}
