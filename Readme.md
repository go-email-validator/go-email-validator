[![codecov](https://codecov.io/gh/go-email-validator/go-email-validator/branch/master/graph/badge.svg?token=BC864E3W3X)](https://codecov.io/gh/go-email-validator/go-email-validator)
[![Go Report](https://goreportcard.com/badge/github.com/go-email-validator/go-email-validator)](https://goreportcard.com/report/github.com/go-email-validator/go-email-validator)

## Install

```go get -u github.com/go-email-validator/go-email-validator```

## Features

username@domain.part

* [syntaxValidator](pkg/ev/validator_syntax.go), using mail.ParseAddress from built-in library
* [disposableValidator](pkg/ev/validator_disposable.go) validation based
  on [mailchecker](https://github.com/FGRibreau/mailchecker)

  Should be fixed, map/hashset can be used instead of array
* [roleValidator](pkg/ev/validator_role.go)
* [mxValidator](pkg/ev/validator_mx.go)
* [smtpValidator](pkg/ev/validator_smtp.go)

    to use proxy set [Checker](pkg/ev/evsmtp/smtp.go) with [proxy_list.ProxySmtpDialer()](pkg/proxifier/proxy_dialer.go)
* [banWordsUsernameValidator](pkg/ev/validator_banwords_username.go), search words in username
* [blackListEmailsValidator](pkg/ev/validator_blacklist_email.go), block emails from list
* [blackListValidator](pkg/ev/validator_blacklist_domain.go), block emails with domain from black list
* [whiteListValidator](pkg/ev/validator_whitelist_domain.go), accept only emails from white list
* [gravatarValidator](pkg/ev/validator_gravatar.go)

## Usage

### With builder

```go
package main

import (
  "fmt"
  "github.com/go-email-validator/go-email-validator/pkg/ev"
  "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

func main() {
  // create defaults list with GetDefaultFactories()
  builder := ev.NewDepBuilder(nil).Build()
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

  v := validator.Validate(evmail.FromString("test@evmail.com"))
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
  "github.com/go-email-validator/go-email-validator/pkg/ev/contains"
  "github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
  "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
  "github.com/go-email-validator/go-email-validator/pkg/ev/free"
  "github.com/go-email-validator/go-email-validator/pkg/ev/role"
  "github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
)

func main() {
  depValidator := ev.NewDepValidator(
    ev.ValidatorMap{
      //ev.FreeValidatorName:       ev.FreeDefaultValidator(),
      ev.RoleValidatorName:       ev.NewRoleValidator(role.NewRBEASetRole()),
      ev.DisposableValidatorName: ev.NewDisposableValidator(contains.NewFunc(disposable.MailChecker)),
      ev.SyntaxValidatorName:     ev.NewSyntaxValidator(),
      ev.MXValidatorName:         ev.NewMXValidator(),
      ev.SMTPValidatorName: ev.NewWarningsDecorator(
        ev.NewSMTPValidator(
          evsmtp.NewChecker(evsmtp.CheckerDTO{
            DialFunc:  evsmtp.Dial,
            SendMail:  evsmtp.NewSendMail(),
            FromEmail: evmail.FromString(evsmtp.DefaultEmail),
          }),
        ),
        ev.NewIsWarning(hashset.New(evsmtp.RandomRCPTStage), func(warningMap ev.WarningSet) ev.IsWarning {
          return func(err error) bool {
            errSMTP, ok := err.(evsmtp.Error)
            if !ok {
              return false
            }
            return warningMap.Contains(errSMTP.Stage())
          }
        }),
      ),
    },
  )

  v := depValidator.Validate(evmail.FromString("test@evmail.com"))
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
  "github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
)

func main() {
  var v = ev.NewSyntaxValidator().Validate(evmail.FromString("some@evmail.here")) // ev.ValidationResult

  if !v.IsValid() {
    panic("email is invalid")
  }

  fmt.Println(v)
}
```

Use function New...(...) to create structure instead of public.

## Addition

1. For running workflow locally use [act](https://github.com/nektos/act)

## Problems

Some mail providers can put ip in spam filter.

1. hotmail.com

## TODO

* Tests
* SMTP working with other ports
* Add linter in pre-hook and ci
* Copy features from [truemail](https://github.com/truemail-rb/truemail)
    * [Extend MX](https://truemail-rb.org/truemail-gem/#/validations-layers?id=mx-validation)
      , [rfc5321 section 5](https://tools.ietf.org/html/rfc5321#section-5)
    * [Host audit features](https://truemail-rb.org/truemail-gem/#/host-audit-features)

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
