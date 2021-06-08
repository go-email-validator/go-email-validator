package checkifemailexist

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

// SyntaxPresentation for check-if-email-exists
type SyntaxPresentation struct {
	Address       *string `json:"address"`
	Username      string  `json:"username"`
	Domain        string  `json:"domain"`
	IsValidSyntax bool    `json:"is_valid_syntax"`
}

// SyntaxConverter converts ev.ValidationResult in SyntaxConverter
type SyntaxConverter struct{}

// Can ev.ValidationResult be converted in SyntaxConverter
func (SyntaxConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.SyntaxValidatorName
}

// Convert ev.ValidationResult in SyntaxPresentation
func (SyntaxConverter) Convert(email evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	presentation := SyntaxPresentation{}

	if result.IsValid() {
		address := email.String()
		presentation.Address = &address
		presentation.Username = email.Username()
		presentation.Domain = email.Domain()
		presentation.IsValidSyntax = result.IsValid()
	}
	return presentation
}
