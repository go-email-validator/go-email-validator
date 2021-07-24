package asemailverifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

// GravatarPresentation is presentation for ev.GravatarValidationResult
type GravatarPresentation struct {
	HasGravatar bool
	GravatarURL string
}

// GravatarConverter converts ev.ValidationResult in GravatarPresentation
type GravatarConverter struct{}

// Can ev.ValidationResult be converted in GravatarPresentation
func (GravatarConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.GravatarValidatorName
}

// Convert ev.ValidationResult in GravatarPresentation
func (GravatarConverter) Convert(_ evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	gravatarResult := result.(ev.GravatarValidationResult)
	presentation := &GravatarPresentation{HasGravatar: gravatarResult.IsValid()}

	if presentation.HasGravatar {
		presentation.GravatarURL = gravatarResult.URL()
	}

	return presentation
}
