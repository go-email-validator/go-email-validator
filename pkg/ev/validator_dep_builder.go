package ev

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/go-email-validator/go-email-validator/pkg/ev/role"
)

type DefaultValidatorFactory func() Validator

func GetDefaultFactories() *ValidatorMap {
	return &ValidatorMap{
		RoleValidatorName:       NewRoleValidator(role.NewRBEASetRole()),
		DisposableValidatorName: NewDisposableValidator(contains.NewFunc(disposable.MailChecker)),
		SyntaxValidatorName:     NewSyntaxValidator(),
		MXValidatorName:         DefaultNewMXValidator(),
		SMTPValidatorName: NewWarningsDecorator(
			smtpValidator{
				checker: evsmtp.NewChecker(evsmtp.CheckerDTO{
					DialFunc:  evsmtp.Dial,
					SendMail:  evsmtp.NewSendMail(),
					FromEmail: evmail.FromString(evsmtp.DefaultEmail),
				}),
			},
			NewIsWarning(hashset.New(evsmtp.RandomRCPTStage), func(warningMap WarningSet) IsWarning {
				return func(err error) bool {
					errSMTP, ok := err.(evsmtp.Error)
					if !ok {
						return false
					}
					return warningMap.Contains(errSMTP.Stage())
				}
			}),
		),
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

func (d *DepBuilder) Has(names ...ValidatorName) bool {
	for _, name := range names {
		if _, has := d.validators[name]; !has {
			return false
		}
	}

	return true
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
