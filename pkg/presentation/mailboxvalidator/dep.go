package mailboxvalidator

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"strconv"
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
	IsFree                EmptyBool     `json:"is_free"`
	IsSyntax              EmptyBool     `json:"is_syntax"`
	IsDomain              EmptyBool     `json:"is_domain"`
	IsSMTP                EmptyBool     `json:"is_smtp"`
	IsVerified            EmptyBool     `json:"is_verified"`
	IsServerDown          EmptyBool     `json:"is_server_down"`
	IsGreylisted          EmptyBool     `json:"is_greylisted"`
	IsDisposable          EmptyBool     `json:"is_disposable"`
	IsSuppressed          EmptyBool     `json:"is_suppressed"`
	IsRole                EmptyBool     `json:"is_role"`
	IsHighRisk            EmptyBool     `json:"is_high_risk"`
	IsCatchall            EmptyBool     `json:"is_catchall"`
	MailboxvalidatorScore float64       `json:"mailboxvalidator_score"`
	TimeTaken             time.Duration `json:"time_taken"`
	Status                EmptyBool     `json:"status"`
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

var depPresentationError = DepPresentation{
	ErrorCode:        strconv.Itoa(int(MissingParameter)),
	CreditsAvailable: ^uint32(0),
}

func getBlankError() DepPresentation {
	return depPresentationError
}

// Convert ev.ValidationResult in mailboxvalidator presentation
func (d DepConverter) Convert(email evmail.Address, resultInterface ev.ValidationResult, opts converter.Options) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			depPresentation := getBlankError()
			depPresentation.ErrorMessage = UnknownErrorMessage

			result = depPresentation
		}
	}()
	var depPresentation DepPresentation
	if len(email.String()) == 0 {
		depPresentation := getBlankError()
		depPresentation.ErrorMessage = MissingParameterMessage

		return depPresentation
	}

	depResult := resultInterface.(ev.DepValidationResult)
	validationResults := depResult.GetResults()

	smtpPresentation := converter.NewSMTPConverter().Convert(email, validationResults[ev.SMTPValidatorName], nil).(converter.SMTPPresentation)

	isFree := !validationResults[ev.FreeValidatorName].IsValid()
	isSyntax := validationResults[ev.SyntaxValidatorName].IsValid()
	depPresentation = DepPresentation{
		EmailAddress:     email.String(),
		Domain:           email.Domain(),
		IsFree:           NewEmptyBool(isFree),
		IsSyntax:         NewEmptyBool(isSyntax),
		IsDomain:         NewEmptyBool(validationResults[ev.MXValidatorName].IsValid()),
		IsSMTP:           NewEmptyBool(smtpPresentation.CanConnectSMTP),
		IsVerified:       NewEmptyBool(smtpPresentation.IsDeliverable),
		IsServerDown:     NewEmptyBool(isSyntax && !smtpPresentation.CanConnectSMTP),
		IsGreylisted:     NewEmptyBool(smtpPresentation.IsGreyListed),
		IsDisposable:     NewEmptyBool(!validationResults[ev.DisposableValidatorName].IsValid()),
		IsSuppressed:     NewEmptyBool(!validationResults[ev.BlackListEmailsValidatorName].IsValid()), // TODO find more examples example@example.com
		IsRole:           NewEmptyBool(!validationResults[ev.RoleValidatorName].IsValid()),
		IsHighRisk:       NewEmptyBool(!validationResults[ev.BanWordsUsernameValidatorName].IsValid()), // TODO find more words
		IsCatchall:       NewEmptyBool(smtpPresentation.IsCatchAll),
		TimeTaken:        opts.ExecutedTime(),
		CreditsAvailable: ^uint32(0),
	}

	depPresentation.MailboxvalidatorScore = d.calculateScore(depPresentation)
	depPresentation.Status = NewEmptyBool(depPresentation.MailboxvalidatorScore >= 0.5)
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
