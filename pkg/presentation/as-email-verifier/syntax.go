package asemailverifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

// SyntaxPresentation is a syntax result
type SyntaxPresentation struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Valid    bool   `json:"valid"`
}

// SyntaxConverter converts ev.ValidationResult in SyntaxPresentation
type SyntaxConverter struct{}

// Can ev.ValidationResult be converted in SyntaxPresentation
func (SyntaxConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.SyntaxValidatorName
}

// Convert ev.ValidationResult in SyntaxPresentation
func (SyntaxConverter) Convert(email evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	presentation := &SyntaxPresentation{Valid: result.IsValid()}

	if presentation.Valid {
		presentation.Username = email.Username()
		presentation.Domain = email.Domain()
	}

	return presentation
}
