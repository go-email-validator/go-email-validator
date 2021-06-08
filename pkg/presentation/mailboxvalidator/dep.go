package mailboxvalidator

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"time"
)

//go:generate go run cmd/dep_test_generator/gen.go

// Name of MailBoxValidator converter
const Name converter.Name = "MailBoxValidator"

// Error constants
const (
	MissingParameter        int32 = 100
	MissingParameterMessage       = "Missing parameter."
	UnknownErrorMessage           = "Unknown error."
)

// DepPresentation is representation of https://www.mailboxvalidator.com/
type DepPresentation struct {
	EmailAddress          string        `json:"email_address"`
	Domain                string        `json:"domain"`
	IsFree                bool          `json:"is_free"`
	IsSyntax              bool          `json:"is_syntax"`
	IsDomain              bool          `json:"is_domain"`
	IsSMTP                bool          `json:"is_smtp"`
	IsVerified            bool          `json:"is_verified"`
	IsServerDown          bool          `json:"is_server_down"`
	IsGreylisted          bool          `json:"is_greylisted"`
	IsDisposable          bool          `json:"is_disposable"`
	IsSuppressed          bool          `json:"is_suppressed"`
	IsRole                bool          `json:"is_role"`
	IsHighRisk            bool          `json:"is_high_risk"`
	IsCatchall            bool          `json:"is_catchall"`
	MailboxvalidatorScore float64       `json:"mailboxvalidator_score"`
	TimeTaken             time.Duration `json:"time_taken"`
	Status                bool          `json:"status"`
	CreditsAvailable      uint32        `json:"credits_available"`
	ErrorCode             string        `json:"error_code"`
	ErrorMessage          string        `json:"error_message"`
}

// FuncCalculateScore is a interface for calculation score function
type FuncCalculateScore func(presentation DepPresentation) float64

// NewDepConverterDefault is the default constructor
func NewDepConverterDefault() DepConverter {
	return NewDepConverter(CalculateScore)
}

// NewDepConverter is the constructor
func NewDepConverter(calculateScore FuncCalculateScore) DepConverter {
	return DepConverter{calculateScore}
}

// DepConverter is the converter for https://www.mailboxvalidator.com/
type DepConverter struct {
	calculateScore FuncCalculateScore
}

// Can ev.ValidationResult be converted in DepConverter
func (DepConverter) Can(_ evmail.Address, result ev.ValidationResult, opts converter.Options) bool {
	return opts.ExecutedTime() != 0 && result.ValidatorName() == ev.DepValidatorName
}

// Convert ev.ValidationResult in mailboxvalidator presentation
func (d DepConverter) Convert(email evmail.Address, resultInterface ev.ValidationResult, opts converter.Options) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			result = DepPresentation{
				ErrorCode:    string(MissingParameter),
				ErrorMessage: UnknownErrorMessage,
			}
		}
	}()
	var depPresentation DepPresentation
	if len(email.String()) == 0 {
		return DepPresentation{
			ErrorCode:    string(MissingParameter),
			ErrorMessage: MissingParameterMessage,
		}
	}

	depResult := resultInterface.(ev.DepValidationResult)
	validationResults := depResult.GetResults()

	smtpPresentation := converter.NewSMTPConverter().Convert(email, validationResults[ev.SMTPValidatorName], nil).(converter.SMTPPresentation)

	isFree := !validationResults[ev.FreeValidatorName].IsValid()
	isSyntax := validationResults[ev.SyntaxValidatorName].IsValid()
	depPresentation = DepPresentation{
		EmailAddress:     email.String(),
		Domain:           email.Domain(),
		IsFree:           isFree,
		IsSyntax:         isSyntax,
		IsDomain:         validationResults[ev.MXValidatorName].IsValid(),
		IsSMTP:           smtpPresentation.CanConnectSMTP,
		IsVerified:       smtpPresentation.IsDeliverable,
		IsServerDown:     isSyntax && !smtpPresentation.CanConnectSMTP,
		IsGreylisted:     smtpPresentation.IsGreyListed,
		IsDisposable:     !validationResults[ev.DisposableValidatorName].IsValid(),
		IsSuppressed:     !validationResults[ev.BlackListEmailsValidatorName].IsValid(), // TODO find more examples example@example.com
		IsRole:           !validationResults[ev.RoleValidatorName].IsValid(),
		IsHighRisk:       !validationResults[ev.BanWordsUsernameValidatorName].IsValid(), // TODO find more words
		IsCatchall:       smtpPresentation.IsCatchAll,
		TimeTaken:        opts.ExecutedTime(),
		CreditsAvailable: ^uint32(0),
	}

	depPresentation.MailboxvalidatorScore = d.calculateScore(depPresentation)
	depPresentation.Status = depPresentation.MailboxvalidatorScore >= 0.5
	return depPresentation
}

// NewDepValidator returns the mailboxvalidator validator
func NewDepValidator(smtpValidator ev.Validator) ev.Validator {
	builder := ev.NewDepBuilder(nil)
	if smtpValidator == nil {
		smtpValidator = builder.Get(ev.SMTPValidatorName)
	}

	return ev.NewDepBuilder(nil).Set(
		ev.BlackListEmailsValidatorName,
		ev.NewBlackListEmailsValidator(contains.NewSet(hashset.New(
			"example@example.com", "localhost@localhost",
		))),
	).Set(
		ev.BanWordsUsernameValidatorName,
		ev.NewBanWordsUsername(contains.NewInStringsFromArray([]string{"test"})),
	).Set(
		ev.FreeValidatorName,
		ev.FreeDefaultValidator(),
	).Set(ev.SMTPValidatorName, smtpValidator).Build()
}
