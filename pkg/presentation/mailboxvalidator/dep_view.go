package mailboxvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"strconv"
	"time"
)

// DepPresentationForView is the DepPresentation but all fields are string
type DepPresentationForView struct {
	EmailAddress          string `json:"email_address"`
	Domain                string `json:"domain"`
	IsFree                string `json:"is_free"`
	IsSyntax              string `json:"is_syntax"`
	IsDomain              string `json:"is_domain"`
	IsSMTP                string `json:"is_smtp"`
	IsVerified            string `json:"is_verified"`
	IsServerDown          string `json:"is_server_down"`
	IsGreylisted          string `json:"is_greylisted"`
	IsDisposable          string `json:"is_disposable"`
	IsSuppressed          string `json:"is_suppressed"`
	IsRole                string `json:"is_role"`
	IsHighRisk            string `json:"is_high_risk"`
	IsCatchall            string `json:"is_catchall"`
	MailboxvalidatorScore string `json:"mailboxvalidator_score"`
	TimeTaken             string `json:"time_taken"`
	Status                string `json:"status"`
	CreditsAvailable      uint32 `json:"credits_available"`
	ErrorCode             string `json:"error_code"`
	ErrorMessage          string `json:"error_message"`
}

// mailbox return time_taken equals 0 for empty email
type jsonAlias DepPresentationForView
type jsonAliasTimeTakenFloat struct {
	TimeTaken        float64 `json:"time_taken"`
	CreditsAvailable string  `json:"credits_available"`
	*jsonAlias
}

func (d *DepPresentationForView) MarshalJSON() ([]byte, error) {
	if d.TimeTaken != "0" {
		return json.Marshal(d)
	}

	var timeTaken float64
	if s, err := strconv.ParseFloat(d.TimeTaken, 32); err == nil {
		timeTaken = s
	}
	aux := jsonAliasTimeTakenFloat{
		TimeTaken: timeTaken,
		jsonAlias: (*jsonAlias)(d),
	}
	return json.Marshal(aux)
}

func (d *DepPresentationForView) UnmarshalJSON(data []byte) error {
	var err error

	aux := (*jsonAlias)(d)
	if err = json.Unmarshal(data, &aux); err == nil {
		return nil
	}

	if errType := err.(*json.UnmarshalTypeError); errType.Field != "time_taken" {
		return err
	}

	auxChanged := jsonAliasTimeTakenFloat{
		jsonAlias: (*jsonAlias)(d),
	}

	if err = json.Unmarshal(data, &auxChanged); err != nil {
		return err
	}
	d.TimeTaken = fmt.Sprint(auxChanged.TimeTaken)
	d.CreditsAvailable = 0

	return nil
}

// FromBool converts bool to string
func FromBool(value bool) string {
	if value {
		return MBVTrue
	}
	return MBVFalse
}

// NewDepConverterForViewDefault creates default DepConverterForView
func NewDepConverterForViewDefault() DepConverterForView {
	return NewDepConverterForView(NewDepConverterDefault())
}

// NewDepConverterForView creates DepConverterForView
func NewDepConverterForView(depConverter DepConverter) DepConverterForView {
	return DepConverterForView{depConverter}
}

// DepConverterForView is the converter for mailbox
type DepConverterForView struct {
	d DepConverter
}

// Can be used to convert ev.ValidationResult in DepConverterForView
func (d DepConverterForView) Can(email evmail.Address, result ev.ValidationResult, opts converter.Options) bool {
	return d.d.Can(email, result, opts)
}

/*
Convert converts the result in mailboxvalidator presentation
TODO add processing of "-" in mailbox validator, for example zxczxczxc@joycasinoru
*/
func (d DepConverterForView) Convert(email evmail.Address, resultInterface ev.ValidationResult, opts converter.Options) interface{} {
	depPresentation := d.d.Convert(email, resultInterface, opts).(DepPresentation)

	return DepPresentationForView{
		EmailAddress:          depPresentation.EmailAddress,
		Domain:                depPresentation.Domain,
		IsFree:                depPresentation.IsFree.ToString(),
		IsSyntax:              depPresentation.IsSyntax.ToString(),
		IsDomain:              depPresentation.IsDomain.ToString(),
		IsSMTP:                depPresentation.IsSMTP.ToString(),
		IsVerified:            depPresentation.IsVerified.ToString(),
		IsServerDown:          depPresentation.IsServerDown.ToString(),
		IsGreylisted:          depPresentation.IsGreylisted.ToString(),
		IsDisposable:          depPresentation.IsDisposable.ToString(),
		IsSuppressed:          depPresentation.IsSuppressed.ToString(),
		IsRole:                depPresentation.IsRole.ToString(),
		IsHighRisk:            depPresentation.IsHighRisk.ToString(),
		IsCatchall:            depPresentation.IsCatchall.ToString(),
		MailboxvalidatorScore: fmt.Sprint(depPresentation.MailboxvalidatorScore),
		TimeTaken:             fmt.Sprint(depPresentation.TimeTaken.Round(time.Microsecond).Seconds()),
		Status:                depPresentation.Status.ToString(),
		CreditsAvailable:      depPresentation.CreditsAvailable,
		ErrorCode:             depPresentation.ErrorCode,
		ErrorMessage:          depPresentation.ErrorMessage,
	}
}
