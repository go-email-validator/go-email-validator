package presenter

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type Reachable string

const (
	Risky   Reachable = "risky"
	Invalid Reachable = "invalid"
	Safe    Reachable = "safe"
	Unknown Reachable = "unknown"
)

func CalculateReachable(depPresenter DepPresenter) Reachable {
	if depPresenter.SMTP == FalseSMTPPresenter {
		if depPresenter.Misc.IsDisposable || depPresenter.Misc.IsRoleAccount || depPresenter.SMTP.IsCatchAll || depPresenter.SMTP.HasFullInbox {
			return Risky
		}

		if !depPresenter.SMTP.IsDeliverable || !depPresenter.SMTP.CanConnectSmtp || depPresenter.SMTP.IsDisabled {
			return Invalid
		}

		return Safe
	} else {
		return Unknown
	}
}

type MiscPresenter struct {
	DisposablePresenter
	RolePresenter
}

type DepPresenter struct {
	Input       string          `json:"input"`
	IsReachable Reachable       `json:"is_reachable"`
	Misc        MiscPresenter   `json:"misc"`
	MX          MXPresenter     `json:"mx"`
	SMTP        SMTPPresenter   `json:"smtp"`
	Syntax      SyntaxPresenter `json:"syntax"`
}

func NewDepProcessor() DepProcessor {
	return DepProcessor{
		MultiplePreparer{MapPreparers{
			ev.RoleValidatorName:       RoleProcessor{},
			ev.DisposableValidatorName: DisposableProcessor{},
			ev.MXValidatorName:         MXProcessor{},
			ev.SMTPValidatorName:       SMTPProcessor{},
			ev.SyntaxValidatorName:     SyntaxProcessor{},
		}},
		CalculateReachable,
	}
}

type DepProcessor struct {
	processor          MultiplePreparer
	calculateReachable func(depPresenter DepPresenter) Reachable
}

func (s DepProcessor) CanProcess(_ email.EmailAddressInterface, result ev.ValidationResultInterface) bool {
	return result.ValidatorName() == ev.DepValidatorName
}

func (s DepProcessor) Process(email email.EmailAddressInterface, result ev.ValidationResultInterface) interface{} {
	presenter := DepPresenter{
		Input: email.String(),
		Misc:  MiscPresenter{},
	}

	for _, result := range result.(ev.DepValidatorResultInterface).GetResults() {
		if !s.processor.CanProcess(email, result) {
			continue
		}

		switch v := s.processor.Process(email, result).(type) {
		case RolePresenter:
			presenter.Misc.RolePresenter = v
		case DisposablePresenter:
			presenter.Misc.DisposablePresenter = v
		case MXPresenter:
			presenter.MX = v
		case SMTPPresenter:
			presenter.SMTP = v
		case SyntaxPresenter:
			presenter.Syntax = v
		}
	}
	presenter.IsReachable = s.calculateReachable(presenter)

	return presenter
}
