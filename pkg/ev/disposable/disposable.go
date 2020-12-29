package disposable

import (
	"github.com/FGRibreau/mailchecker/v4/platform/go"
)

// Send domain to check.
func MailChecker(domain interface{}) bool {
	return mail_checker.IsBlacklisted("username@" + domain.(string))
}
