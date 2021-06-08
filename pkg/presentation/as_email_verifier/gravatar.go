package as_email_verifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

type GravatarPresentation struct {
	HasGravatar bool
	GravatarUrl string
}

type GravatarConverter struct{}

func (GravatarConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.GravatarValidatorName
}

func (GravatarConverter) Convert(_ evmail.Address, result ev.ValidationResult, _ converter.Options) interface{} {
	gravatarResult := result.(ev.GravatarValidationResult)
	presentation := &GravatarPresentation{HasGravatar: gravatarResult.IsValid()}

	if presentation.HasGravatar {
		presentation.GravatarUrl = gravatarResult.URL()
	}

	return presentation
}
