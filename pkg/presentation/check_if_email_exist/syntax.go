package check_if_email_exist

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

type syntaxPresentation struct {
	Address       *string `json:"address"`
	Username      string  `json:"username"`
	Domain        string  `json:"domain"`
	IsValidSyntax bool    `json:"is_valid_syntax"`
}

type SyntaxConverter struct{}

func (SyntaxConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.SyntaxValidatorName
}

func (SyntaxConverter) Convert(email evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	presentation := syntaxPresentation{}

	if result.IsValid() {
		address := email.String()
		presentation.Address = &address
		presentation.Username = email.Username()
		presentation.Domain = email.Domain()
		presentation.IsValidSyntax = result.IsValid()
	}
	return presentation
}
