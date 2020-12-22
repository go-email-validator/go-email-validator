package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
	"github.com/go-email-validator/go-email-validator/pkg/ev/smtp_checker"
)

type DefaultValidatorFactory func() ValidatorInterface

func GetDefaultFactories() *ValidatorMap {
	return &ValidatorMap{
		DisposableValidatorName: NewDisposableValidator(contains.NewFunc(disposable.MailChecker)),
		RoleValidatorName:       NewRoleValidator(role.NewRBEASetRole()),
		MXValidatorName:         NewMXValidator(),
		SMTPValidatorName: NewWarningsDecorator(
			SMTPValidator{
				Checker: smtp_checker.Checker{
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

func NewDepBuilder(validators *ValidatorMap) DepBuilder {
	if validators == nil {
		validators = GetDefaultFactories()
	}

	return DepBuilder{validators: *validators}
}

type DepBuilder struct {
	validators ValidatorMap
}

func (d DepBuilder) Set(name ValidatorName, validator ValidatorInterface) {
	d.validators[name] = validator
}

func (d DepBuilder) Has(name ValidatorName) bool {
	_, has := d.validators[name]

	return has
}

func (d DepBuilder) Delete(name ValidatorName) {
	if d.Has(name) {
		delete(d.validators, name)
	}
}

func (d DepBuilder) Build() ValidatorInterface {
	return NewDepValidator(d.validators)
}
