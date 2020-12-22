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

func NewRoleValidator(r contains.Interface) ValidatorInterface {
	return RoleValidator{r: r}
}

type RoleValidator struct {
	r contains.Interface
	AValidatorWithoutDeps
}

func (r RoleValidator) Validate(email ev_email.EmailAddressInterface, _ ...ValidationResultInterface) ValidationResultInterface {
	var err error
	var hasRole = r.r.Contains(email.Username())
	if hasRole {
		err = RoleError{}
	}

	return NewValidatorResult(!hasRole, utils.Errs(err), nil, RoleValidatorName)
}
