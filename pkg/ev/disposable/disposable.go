package disposable

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
	mail_checker "github.com/FGRibreau/mailchecker/platform/go"
)

type Interface interface {
	Disposable(email ev_email.EmailAddressInterface) bool
}

type SetDisposable struct {
	set utils.StringSet
}

func (s SetDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	_, ok := s.set[email.Domain()]
	return ok
}

type MailCheckerDisposable struct{}

func (m MailCheckerDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return mail_checker.IsBlacklisted(email.String())
}
