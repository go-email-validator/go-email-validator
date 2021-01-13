package disposable

import (
	"github.com/FGRibreau/mailchecker/v4/platform/go"
)

// MailChecker sends domain to check by https://github.com/FGRibreau/mailchecker/
func MailChecker(domain interface{}) bool {
	return mail_checker.IsBlacklisted("username@" + domain.(string))
}
