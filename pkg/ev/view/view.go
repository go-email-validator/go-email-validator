package view

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev"
	email "bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
)

type ViewInterface interface {
}

type FactoryView = func(email email.EmailAddressInterface, result ev.ValidationResultInterface) ViewInterface
