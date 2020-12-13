## How to use

```go
depValidator := DepValidator{
    map[ValidatorName]ValidatorInterface{
        RoleValidatorName: NewRoleValidator(role.NewRBEASetRole()),
        DisposableValidatorName: NewDisposableValidator(disposable.MailCheckerDisposable{}),
        SyntaxValidatorName: &SyntaxValidator{},
        MXValidatorName:     &MXValidator{},
        SMTPValidatorName: NewWarningsDecorator(
            ValidatorInterface(newSMTPValidator()),
            NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
                return func(err error) bool {
                    return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
                }
            }),
        ),
    },
}

v := depValidator.Validate(email)

if !v.isValid() {
    panic('email is invalid')
}
```

## TODO

* Builder for [DepValidator](pkg/ev/validator_dep.go)
* Tests

## Inspired by

* [EmailValidator](https://github.com/egulias/EmailValidator)
