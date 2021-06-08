package prompt_email_verification_api

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/presentation/converter"
	"strings"
)

var emptyString = ""

func EmailFromString(email string) evmail.Address {
	firstPos := strings.IndexByte(email, '@')
	lastPos := strings.LastIndexByte(email, '@')

	if firstPos == -1 || len(email) < 3 || firstPos != lastPos {
		return converter.NewEmailAddress("", "", &emptyString)
	}

	return converter.NewEmailAddress(email[:firstPos], email[firstPos+1:], nil)
}
