package mailboxvalidator

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"strings"
)

var emptyString = ""

// EmailFromString creates evmail.Address from string
func EmailFromString(email string) evmail.Address {
	pos := strings.LastIndexByte(email, '@')

	if pos == -1 || len(email) < 3 {
		return converter.NewEmailAddress("", email, &emptyString)
	}

	return converter.NewEmailAddress(email[:pos], email[pos+1:], nil)
}
