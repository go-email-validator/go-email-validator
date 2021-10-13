[![Go Reference](https://pkg.go.dev/badge/github.com/go-email-validator/go-email-validator.svg)](https://pkg.go.dev/github.com/go-email-validator/go-email-validator)
[![codecov](https://codecov.io/gh/go-email-validator/go-email-validator/branch/master/graph/badge.svg?token=BC864E3W3X)](https://codecov.io/gh/go-email-validator/go-email-validator)
[![Go Report](https://goreportcard.com/badge/github.com/go-email-validator/go-email-validator)](https://goreportcard.com/report/github.com/go-email-validator/go-email-validator)

## Library under development (Interfaces may be changed slightly)

## Demo on [rapidapi.com](https://rapidapi.com/maranqz/api/email-validator15)

## Install

```go get -u github.com/go-email-validator/go-email-validator```

## Available validators

username@domain.part

* [syntaxValidator](pkg/ev/validator_syntax.go)
  1. `NewSyntaxValidator()` - mail.ParseAddress from built-in library
  1. `NewSyntaxRegexValidator(emailRegex *regexp.Regexp)` - validation based on regular expression
* [disposableValidator](pkg/ev/validator_disposable.go) based
  on [mailchecker](https://github.com/FGRibreau/mailchecker) by default (set is replaceable)
* [roleValidator](pkg/ev/validator_role.go) bases on [role-based-email-addresses](https://github.com/mixmaxhq/role-based-email-addresses) by default (set is replaceable)
* [mxValidator](pkg/ev/validator_mx.go)
* [smtpValidator](pkg/ev/validator_smtp.go)

    to use proxy connection, DialFunc need to be changed in [Checker](pkg/ev/evsmtp/smtp.go). There is [evsmtp.H12IODial](pkg/ev/evsmtp/proxy.go), implementing for [h12w](https://github.com/h12w/socks).
* [banWordsUsernameValidator](pkg/ev/validator_banwords_username.go) looks for banned words in username
* [blackListEmailsValidator](pkg/ev/validator_blacklist_email.go) blocked emails from list
* [blackListValidator](pkg/ev/validator_blacklist_domain.go) blocked emails with domains from black list
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

  v := validator.Validate(NewInput(evmail.FromString("test@evmail.com")))
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
  var v = ev.NewSyntaxValidator().Validate(NewInput(evmail.FromString("some@evmail.here"))) // ev.ValidationResult

  if !v.IsValid() {
    panic("email is invalid")
  }

  fmt.Println(v)
}
```

### Addition options

To set options for different validators, use NewInput(..., NewKVOption(ValidatorName, Options))

```go
NewInput(
	evmail.FromString("test@evmail.com"), 
	NewKVOption(SMTPValidatorName, evsmtp.NewOptions(evsmtp.OptionsDTO{
		Port: 465,
    })),
)
```

Use function New...(...) to create structure instead of public.

## How to extend

To add own validator, just implement [ev.Validator](pkg/ev/validator.go) interface. For validator without dependencies, you can use structure ev.AValidatorWithoutDeps

## Decorators

1. [WarningsDecorator](pkg/ev/decorator_warnings.go) allows moving errors to warnings and change result of `IsValid()` in [ValidationResult](pkg/ev/validator.go).
1. Cache based on evcahce.Interface, default realization is done for gocache.
    * [CacheDecorator](pkg/ev/decorator_cache.go) saves result of validator. For caching, you can implement `evcache.Interface` or use [gocache implementation](https://github.com/eko/gocache) by `evcache.NewCache`. See [Test_Cache](pkg/ev/decorator_cache_test.go) as example.
    * [checkerCacheRandomRCPT](pkg/ev/evsmtp/smtp.go) for caching of RandomRCPTs request. See Test_checkerCacheRandomRCPT_RandomRCPT_RealCache as example.

**Notice**, to use [msgpack](https://github.com/vmihailenco/msgpack) you should have exported fields or implement custom encoding/decoding ([doc](https://msgpack.uptrace.dev/#custom-encodingdecoding))

## Logger

Package use [zap](https://github.com/uber-go/zap).

To use logging see in [log package](pkg/log).
Default level is zap.ErrorLevel.

## Addition

1. For running workflow locally use [act](https://github.com/nektos/act)

## FAQ

#### Most Internet Service Providers block outgoing SMTP request.

The [StackOverflow thread](https://stackoverflow.com/questions/18139102/how-to-get-around-an-isp-block-on-port-25-for-smtp) could be helpful.

#### To check smtp in telnet

```
telnet

OPEN gmail-smtp-in.l.google.com 25
EHLO localhost
MAIL FROM: <user@example.org>
rcpt to: <some.email@gmail.com>
quit
```

#### Some mail providers could put your ip in spam filter.

For example:
1. hotmail.com

## Roadmap

* Tests
  * Add functional tests
  * Find way to compare functions in tests
* Add binary release
* Check in spamhaus
* Add misspelled email
* Add DKIM checking
* Add linter in pre-hook and ci
* Do full thread safe library
* Copy features from [truemail](https://github.com/truemail-rb/truemail)
    * [Extend MX](https://truemail-rb.org/truemail-gem/#/validations-layers?id=mx-validation)
      , [rfc5321 section 5](https://tools.ietf.org/html/rfc5321#section-5)
    * [Host audit features](https://truemail-rb.org/truemail-gem/#/host-audit-features)

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
