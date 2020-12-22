package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
)

type DefaultValidatorFactory func() Validator

func GetDefaultFactories() *ValidatorMap {
	return &ValidatorMap{
		DisposableValidatorName: NewDisposableValidator(contains.NewFunc(disposable.MailChecker)),
		RoleValidatorName:       NewRoleValidator(role.NewRBEASetRole()),
		MXValidatorName:         NewMXValidator(),
		SMTPValidatorName: NewWarningsDecorator(
			smtpValidator{
				checker: smtp_checker.Checker{
					GetConn:   smtp_checker.SimpleClientGetter,
					SendMail:  smtp_checker.NewSendMail(),
					FromEmail: ev_email.EmailFromString(smtp_checker.DefaultEmail),
				},
			},
			NewIsWarning(hashset.New(smtp_checker.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
				return func(err error) bool {
					return warningMap.Contains(err.(smtp_checker.SMTPError).Stage())
				}
			}),
		),
		SyntaxValidatorName: NewSyntaxValidator(),
	}
}

func NewDepBuilder(validators *ValidatorMap) *DepBuilder {
	if validators == nil {
		validators = GetDefaultFactories()
	}

	return &DepBuilder{validators: *validators}
}

type DepBuilder struct {
	validators ValidatorMap
}

func (d *DepBuilder) Set(name ValidatorName, validator Validator) *DepBuilder {
	d.validators[name] = validator

	return d
}

func (d *DepBuilder) Has(name ValidatorName) bool {
	_, has := d.validators[name]

	return has
}

func (d *DepBuilder) Delete(names ...ValidatorName) *DepBuilder {
	for _, name := range names {
		if d.Has(name) {
			delete(d.validators, name)
		}
	}

	return d
}

func (d *DepBuilder) Build() Validator {
	return NewDepValidator(d.validators)
}
