## Features

username@domain.part

* [syntaxValidator](pkg/ev/validator_syntax.go), using mail.ParseAddress from built-in library
* [disposableValidator](pkg/ev/validator_disposable.go) validation based
  on [mailchecker](https://github.com/FGRibreau/mailchecker)

  Should be fixed, map/hashset can be used instead of array
* [roleValidator](pkg/ev/validator_role.go)
* [mxValidator](pkg/ev/validator_mx.go)
* [smtpValidator](pkg/ev/validator_smtp.go)

    to use proxy set [Checker](pkg/ev/smtp_checker/smtp.go) with [proxy_list.ProxySmtpDialer()](pkg/proxy_list/proxy_dialer.go)
* [banWordsUsernameValidator](pkg/ev/validator_banwords_username.go), search words in username
* [blackListEmailsValidator](pkg/ev/validator_blacklist_email.go), block emails from list
* [blackListValidator](pkg/ev/validator_blacklist_domain.go), block emails with domain from black list
* [whiteListValidator](pkg/ev/validator_whitelist_domain.go), accept only emails from white list
* [gravatarValidator](pkg/ev/validator_gravatar.go)

## How to use

### With builder

```go
package main

import (
  "fmt"
  "github.com/go-email-validator/go-email-validator/pkg/ev"
)

func main() {
  // create defaults list with GetDefaultFactories()
  builder := NewDepBuilder().Build()
  /*
     to set list of initial validators
     builder := NewDepBuilder(&ValidatorMap{
         key: ev.Validator,
     }).Build()
  */

  // validator.Set(NameValidator, NewValidator())  builder
  // validator.Has(names ...ValidatorName) bool
  // validator.Delete(names ...ValidatorName) bool

  validator := builder.Build()

  v := validator.Validate(ev_email.EmailFromString("test@email.com"))
  if !v.IsValid() {
    panic("email is invalid")
  }

  fmt.Println(v)
}

```

### Clean

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
  "net/smtp"
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
          smtp_checker.NewChecker(smtp_checker.CheckerDTO{
            DialFunc:  smtp.Dial,
            SendMail:  smtp_checker.NewSendMail(),
            FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
          }),
        ),
        ev.NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap ev.WarningSet) ev.IsWarning {
          return func(err error) bool {
            errSmtp, ok := err.(smtp_checker.SMTPError)
            if !ok {
              return false
            }
            return warningMap.Contains(errSmtp.Stage())
          }
        }),
      ),
    },
  )

  v := depValidator.Validate(ev_email.EmailFromString("test@email.com"))
  if !v.IsValid() {
    panic("email is invalid")
  }

  fmt.Println(v)
}
```

### Single validator

```go
package main

import (
  "fmt"
  "github.com/go-email-validator/go-email-validator/pkg/ev"
  "github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
)

func main() {
  var v = ev.NewSyntaxValidator().Validate(ev_email.EmailFromString("some@email.here")) // ev.ValidationResult

  if !v.IsValid() {
    panic("email is invalid")
  }

  fmt.Println(v)
}
```

Use function New...(...) to create structure instead of public.

## TODO

* Tests
* SMTP working with other ports
* Copy features from [truemail](https://github.com/truemail-rb/truemail)
    * [Extend MX](https://truemail-rb.org/truemail-gem/#/validations-layers?id=mx-validation)
      , [rfc5321 section 5](https://tools.ietf.org/html/rfc5321#section-5)
    * [Host audit features](https://truemail-rb.org/truemail-gem/#/host-audit-features)

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
