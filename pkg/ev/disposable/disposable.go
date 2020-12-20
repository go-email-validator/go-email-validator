package disposable

import (
	mail_checker "github.com/FGRibreau/mailchecker/v4/platform/go"
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	Disposable(email ev_email.EmailAddressInterface) bool
}

type SetDisposable struct {
	set sets.Set
}

func (s SetDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Domain())
}

// List is used for searching blacklisted email domain
type MailCheckerDisposable struct{}

func (_ MailCheckerDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return mail_checker.IsBlacklisted(email.String())
}
