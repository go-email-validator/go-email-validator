## How to use

```go
package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
)

func main() {
	depValidator := ev.NewDepValidator(
		map[ev.ValidatorName]ev.ValidatorInterface{
			//ev.FreeValidatorName:       ev.FreeDefaultValidator(),
			ev.RoleValidatorName:       ev.NewRoleValidator(role.NewRBEASetRole()),
			ev.DisposableValidatorName: ev.NewDisposableValidator(disposable.NewFuncDisposable(disposable.MailChecker)),
			ev.SyntaxValidatorName:     ev.NewSyntaxValidator(),
			ev.MXValidatorName:         ev.NewMXValidator(),
			ev.SMTPValidatorName: ev.NewWarningsDecorator(
				ev.NewSMTPValidator(
					smtp_checker.Checker{
						GetConn:   smtp_checker.SimpleClientGetter,
						SendMail:  smtp_checker.NewSendMail(),
						FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
					},
                ),
				ev.NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap ev.WarningSet) ev.IsWarning {
					return func(err error) bool {
						return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
					}
				}),
			),
		},
	)

	v := depValidator.Validate(ev_email.EmailFromString("test@email.com"))
	if !v.isValid() {
		panic('email is invalid')
	}

	fmt.Print(v)
}
```

Use func New...(...) instead of public struct.

## TODO

* Builder for [DepValidator](pkg/ev/validator_dep.go)
* Tests
* Copy features from [truemail](https://github.com/truemail-rb/truemail)
    * [Extend MX](https://truemail-rb.org/truemail-gem/#/validations-layers?id=mx-validation)
      , [rfc5321 section 5](https://tools.ietf.org/html/rfc5321#section-5)
    * [Host audit features](https://truemail-rb.org/truemail-gem/#/host-audit-features)

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
