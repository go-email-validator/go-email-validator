package asemailverifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"regexp"
)

const (
	// Name of https://github.com/AfterShip/email-verifier converter
	Name converter.Name = "AfterShipEmailVerifier"

	// EmailRegexString is an email regex
	EmailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

// DepPresentation is a presentation
type DepPresentation struct {
	Email        string                `json:"email"`
	Reachable    Reachable             `json:"reachable"`
	Syntax       *SyntaxPresentation   `json:"syntax"`
	SMTP         *SMTPPresentation     `json:"smtp"`
	Gravatar     *GravatarPresentation `json:"gravatar"`
	Suggestion   string                `json:"suggestion"`
	Disposable   bool                  `json:"disposable"`
	RoleAccount  bool                  `json:"role_account"`
	Free         bool                  `json:"free"`
	HasMxRecords bool                  `json:"has_mx_records"`
}

// DepConverter converts ev.ValidationResult in DepPresentation
type DepConverter struct {
	converter          converter.CompositeConverter
	calculateReachable FuncReachable
}

// NewDepConverterDefault creates default DepConverter
func NewDepConverterDefault() DepConverter {
	return NewDepConverter(
		converter.NewCompositeConverter(converter.MapConverters{
			ev.SMTPValidatorName:     converter.NewSMTPConverter(),
			ev.SyntaxValidatorName:   SyntaxConverter{},
			ev.GravatarValidatorName: GravatarConverter{},
		}),
		CalculateReachable,
	)
}

// NewDepConverter creates DepConverter
func NewDepConverter(converter converter.CompositeConverter, calculateReachable FuncReachable) DepConverter {
	return DepConverter{converter, calculateReachable}
}

// Can ev.ValidationResult be converted in DepConverter
func (DepConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.DepValidatorName
}

// Convert ev.ValidationResult in DepConverter
func (s DepConverter) Convert(email evmail.Address, resultInterface ev.ValidationResult, opts converter.Options) interface{} {
	depResult := resultInterface.(ev.DepValidationResult)
	validationResults := depResult.GetResults()

	var syntax *SyntaxPresentation
	var smtp *SMTPPresentation
	var gravatar *GravatarPresentation
	mxResult := validationResults[ev.MXValidatorName].(ev.MXValidationResult)
	for _, validatorResult := range depResult.GetResults() {
		if !s.converter.Can(email, validatorResult, opts) {
			continue
		}

		switch v := s.converter.Convert(email, validatorResult, opts).(type) {
		case *SyntaxPresentation:
			syntax = v
		case converter.SMTPPresentation:
			smtp = &SMTPPresentation{
				HostExists:  v.CanConnectSMTP,
				FullInbox:   v.HasFullInbox,
				CatchAll:    v.IsCatchAll,
				Deliverable: !v.IsCatchAll && v.IsDeliverable,
				Disabled:    v.IsDisabled,
			}
		case *GravatarPresentation:
			gravatar = v
		}
	}

	depPresentation := DepPresentation{
		Email:     email.String(),
		Reachable: ReachableUnknown,
		Syntax:    syntax,
	}

	if syntax == nil || !syntax.Valid {
		return depPresentation
	}

	depPresentation.Free = !validationResults[ev.FreeValidatorName].IsValid()
	depPresentation.RoleAccount = !validationResults[ev.RoleValidatorName].IsValid()
	depPresentation.Disposable = !validationResults[ev.DisposableValidatorName].IsValid()

	if depPresentation.Disposable {
		return depPresentation
	}

	if !mxResult.IsValid() {
		return depPresentation
	}

	depPresentation.HasMxRecords = len(mxResult.MX()) > 0

	if smtp == nil || !smtp.HostExists {
		return depPresentation
	}

	depPresentation.SMTP = smtp
	depPresentation.Reachable = s.calculateReachable(depPresentation)

	if gravatar != nil {
		depPresentation.Gravatar = gravatar
	}

	return depPresentation
}

// NewDepValidator is a default validator
func NewDepValidator(smtpValidator ev.Validator) ev.Validator {
	builder := ev.NewDepBuilder(nil)
	if smtpValidator == nil {
		smtpValidator = builder.Get(ev.SMTPValidatorName)
	}

	return ev.NewDepBuilder(nil).Set(
		ev.SyntaxValidatorName,
		ev.NewSyntaxRegexValidator(regexp.MustCompile(EmailRegexString)),
	).
		Set(ev.GravatarValidatorName, ev.NewGravatarValidator()).
		Set(ev.SMTPValidatorName, smtpValidator).
		Set(ev.FreeValidatorName, ev.FreeDefaultValidator()).
		Build()
}
