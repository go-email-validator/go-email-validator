package disposable

import (
	mail_checker "github.com/FGRibreau/mailchecker/v4/platform/go"
	"github.com/emirpasic/gods/sets"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

type Interface interface {
	Disposable(email ev_email.EmailAddressInterface) bool
}

func NewSetDisposable(s sets.Set) Interface {
	return setDisposable{s}
}

type setDisposable struct {
	set sets.Set
}

func (s setDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return s.set.Contains(email.Domain())
}

type FuncChecker func(email ev_email.EmailAddressInterface) bool

func NewFuncDisposable(f FuncChecker) Interface {
	return funcDisposable{f}
}

type funcDisposable struct {
	funcChecker FuncChecker
}

func (f funcDisposable) Disposable(email ev_email.EmailAddressInterface) bool {
	return f.funcChecker(email)
}

// List is used for searching blacklisted email domain
func MailChecker(email ev_email.EmailAddressInterface) bool {
	return mail_checker.IsBlacklisted(email.String())
}
