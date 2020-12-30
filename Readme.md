[![Go Reference](https://pkg.go.dev/badge/github.com/go-email-validator/go-email-validator.svg)](https://pkg.go.dev/github.com/go-email-validator/go-email-validator)
[![codecov](https://codecov.io/gh/go-email-validator/go-email-validator/branch/master/graph/badge.svg?token=BC864E3W3X)](https://codecov.io/gh/go-email-validator/go-email-validator)
[![Go Report](https://goreportcard.com/badge/github.com/go-email-validator/go-email-validator)](https://goreportcard.com/report/github.com/go-email-validator/go-email-validator)

## Library under development (Interfaces may change slightly)

## Install

```go get -u github.com/go-email-validator/go-email-validator```

## Available validators

username@domain.part

* [syntaxValidator](pkg/ev/validator_syntax.go) uses mail.ParseAddress from built-in library
* [disposableValidator](pkg/ev/validator_disposable.go) based
  on [mailchecker](https://github.com/FGRibreau/mailchecker) by default (set is replaceable)
* [roleValidator](pkg/ev/validator_role.go) bases on [role-based-email-addresses](https://github.com/mixmaxhq/role-based-email-addresses) by default (set is replaceable)
* [mxValidator](pkg/ev/validator_mx.go)
* [smtpValidator](pkg/ev/validator_smtp.go)

    to use proxy connection DialFunc need to be changed for [Checker](pkg/ev/evsmtp/smtp.go). For example by [ProxyDialer](pkg/proxifier/proxy_dialer.go)
* [banWordsUsernameValidator](pkg/ev/validator_banwords_username.go) looks for banned words in username
* [blackListEmailsValidator](pkg/ev/validator_blacklist_email.go) blocks emails from list
* [blackListValidator](pkg/ev/validator_blacklist_domain.go) blocks emails with domain from black list
* [whiteListValidator](pkg/ev/validator_whitelist_domain.go) accepts only emails from white list
* [gravatarValidator](pkg/ev/validator_gravatar.go) check existing of user on gravatar.com

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
  // create defaults DepValidator with GetDefaultFactories() as list of validators
  builder := ev.NewDepBuilder(nil).Build()
  /*
     to set another list of initial validators
     builder := NewDepBuilder(&ValidatorMap{
         ev.ValidatorName: ev.Validator,
     }).Build()
  */

  // builder.Set(ev.ValidatorName, NewValidator()) builder
  // builder.Has(names ...ev.ValidatorName) bool
  // builder.Delete(names ...ev.ValidatorName) bool

  validator := builder.Build()

  v := validator.Validate(evmail.FromString("test@evmail.com"))
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

## How to extend

To add own validator, just realize [ev.Validator](pkg/ev/validator.go) interface. For validator without dependencies, you can use structure ev.AValidatorWithoutDeps

## Addition

1. For running workflow locally use [act](https://github.com/nektos/act)

## FAQ

#### Most Internet Service Providers block outgoing SMTP request.

The [StackOverflow thread](https://stackoverflow.com/questions/18139102/how-to-get-around-an-isp-block-on-port-25-for-smtp) could be helpful.

#### Some mail providers could put your ip in spam filter.

For example:
1. hotmail.com

## Roadmap

* Tests
  * Add functional tests
  * Find way to compare functions in tests
* Add the ability to work SMTP with other ports
* Add regexp syntax validation
* Add linter in pre-hook and ci
* Copy features from [truemail](https://github.com/truemail-rb/truemail)
    * [Extend MX](https://truemail-rb.org/truemail-gem/#/validations-layers?id=mx-validation)
      , [rfc5321 section 5](https://tools.ietf.org/html/rfc5321#section-5)
    * [Host audit features](https://truemail-rb.org/truemail-gem/#/host-audit-features)

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
