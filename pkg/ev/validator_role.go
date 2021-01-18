package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
)

// RoleValidatorName is name of role validator
const RoleValidatorName ValidatorName = "RoleValidator"

// RoleError is error of RoleValidatorName
type RoleError struct{}

func (RoleError) Error() string {
	return "RoleError"
}

// NewRoleValidator instantiates RoleValidatorName
func NewRoleValidator(r contains.InSet) Validator {
	return roleValidator{r: r}
}

type roleValidator struct {
	AValidatorWithoutDeps
	r contains.InSet
}

func (r roleValidator) Validate(input Input, _ ...ValidationResult) ValidationResult {
	var err error
	var hasRole = r.r.Contains(input.Email().Username())
	if hasRole {
		err = RoleError{}
	}

	return NewResult(!hasRole, utils.Errs(err), nil, RoleValidatorName)
}
