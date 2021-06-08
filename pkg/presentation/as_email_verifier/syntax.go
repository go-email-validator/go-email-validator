package as_email_verifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

type SyntaxPresentation struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Valid    bool   `json:"valid"`
}

type SyntaxConverter struct{}

func (SyntaxConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.SyntaxValidatorName
}

func (SyntaxConverter) Convert(email evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	presentation := &SyntaxPresentation{Valid: result.IsValid()}

	if presentation.Valid {
		presentation.Username = email.Username()
		presentation.Domain = email.Domain()
	}

	return presentation
}
