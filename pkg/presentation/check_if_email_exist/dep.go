package check_if_email_exist

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
)

// TODO create transfer from DepPresentationForView to DepPresentation
//go:generate go run cmd/dep_test_generator/gen.go

const Name converter.Name = "CheckIfEmailExist"

type miscPresentation struct {
	disposablePresentation
	rolePresentation
}

// DepPresentation is representation of https://github.com/amaurymartiny/check-if-email-exists or https://reacher.email
type DepPresentation struct {
	Input       string             `json:"input"`
	IsReachable Availability       `json:"is_reachable"`
	Misc        miscPresentation   `json:"misc"`
	MX          mxPresentation     `json:"mx"`
	SMTP        SmtpPresentation   `json:"smtp"`
	Syntax      syntaxPresentation `json:"syntax"`
	Error       string             `json:"error"`
}

type DepConverter struct {
	converter             converter.CompositeConverter
	calculateAvailability FuncAvailability
}

func NewDepConverterDefault() DepConverter {
	return NewDepConverter(
		converter.NewCompositeConverter(converter.MapConverters{
			ev.RoleValidatorName:       roleConverter{},
			ev.DisposableValidatorName: disposableConverter{},
			ev.MXValidatorName:         mxConverter{},
			ev.SMTPValidatorName:       converter.NewSMTPConverter(),
			ev.SyntaxValidatorName:     SyntaxConverter{},
		}),
		CalculateAvailability,
	)
}

func NewDepConverter(converter converter.CompositeConverter, calculateAvailability FuncAvailability) DepConverter {
	return DepConverter{converter, calculateAvailability}
}

func (DepConverter) Can(_ evmail.Address, result ev.ValidationResult, _ converter.Options) bool {
	return result.ValidatorName() == ev.DepValidatorName
}

func (s DepConverter) Convert(email evmail.Address, result ev.ValidationResult, opts converter.Options) interface{} {
	depPresentation := DepPresentation{
		Input: email.String(),
		Misc:  miscPresentation{},
	}

	for _, validatorResult := range result.(ev.DepValidationResult).GetResults() {
		if !s.converter.Can(email, validatorResult, opts) {
			continue
		}

		switch v := s.converter.Convert(email, validatorResult, opts).(type) {
		case rolePresentation:
			depPresentation.Misc.rolePresentation = v
		case disposablePresentation:
			depPresentation.Misc.disposablePresentation = v
		case mxPresentation:
			depPresentation.MX = v
		case converter.SmtpPresentation:
			depPresentation.SMTP = SmtpPresentation{
				CanConnectSmtp: v.CanConnectSmtp,
				HasFullInbox:   v.HasFullInbox,
				IsCatchAll:     v.IsCatchAll,
				IsDeliverable:  v.IsDeliverable,
				IsDisabled:     v.IsDisabled,
			}
		case syntaxPresentation:
			depPresentation.Syntax = v
		}
	}
	depPresentation.IsReachable = s.calculateAvailability(depPresentation)

	return depPresentation
}

func NewDepValidator(smtpValidator ev.Validator) ev.Validator {
	builder := ev.NewDepBuilder(nil)
	if smtpValidator == nil {
		smtpValidator = builder.Get(ev.SMTPValidatorName)
	}

	return builder.Set(
		ev.SyntaxValidatorName,
		ev.NewSyntaxRegexValidator(nil),
	).Set(
		ev.SMTPValidatorName,
		smtpValidator,
	).Build()
}
