package disposable

import (
	"github.com/FGRibreau/mailchecker/v4/platform/go"
)

// List is used for searching blacklisted email domain
func MailChecker(value interface{}) bool {
	return mail_checker.IsBlacklisted(value.(string))
}
