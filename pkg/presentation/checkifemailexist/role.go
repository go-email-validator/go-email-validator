package checkifemailexist

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

type rolePresentation struct {
	IsRoleAccount bool `json:"is_role_account"`
}

type roleConverter struct{}

func (roleConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.RoleValidatorName
}

func (roleConverter) Convert(_ evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	return rolePresentation{!result.IsValid()}
}
