package disposable

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	mail_checker "github.com/FGRibreau/mailchecker/platform/go"
)

type Void struct{}

func GetVoid() Void {
	var member Void

	return member
}

type StringSet map[string]Void

type Interface interface {
	Disposable(email ev_email.EmailAddressInterface) bool
}

type SetDisposable struct {
	set StringSet
}

func (s SetDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	_, ok := s.set[email.GetDomain()]
	return ok
}

type MailCheckerDisposable struct{}

func (m MailCheckerDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return mail_checker.IsBlacklisted(email.String())
}
