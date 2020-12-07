package smtp_checker

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"bitbucket.org/maranqz/email-validator/pkg/ev/utils"
	"errors"
	"fmt"
	"github.com/sethvargo/go-password/password"
	"net/smtp"
)

const (
	SMTPErrorHelloAfter = "smtp_checker: Hello called after other methods"
	SMTPErrorCrLR       = "smtp_checker: A line must not contain CR or LF"

	DefaultEmail = "user@example.org"

	DefaultSMTPPort = 25
)

type SMTPError interface {
	error
	Stage() SendMailStage
	Err() error
}

type ASMTPError struct {
	stage SendMailStage
	err   error
}

func (a ASMTPError) Stage() SendMailStage {
	return a.stage
}
func (a ASMTPError) Err() error {
	return a.err
}

func (a ASMTPError) Error() string {
	return fmt.Sprintf("%v happend on stage \"%v\"", a.Err().Error(), a.Stage())
}

func NewSmtpError(stage SendMailStage, err error) SMTPError {
	return DefaultSmtpError{ASMTPError{stage, err}}
}

type DefaultSmtpError struct {
	ASMTPError
}

type SMTPErrorNested interface {
	SMTPError
	GetNested() SMTPError
}

type ASMTPErrorNested struct {
	n SMTPError
}

func (a ASMTPErrorNested) GetNested() SMTPError {
	return a.n
}

func (a ASMTPErrorNested) Error() string {
	return a.n.Error()
}

func randomEmail(domain string) (ev_email.EmailAddressInterface, error) {
	input := new(password.GeneratorInput)
	input.LowerLetters = password.LowerLetters + password.Digits

	gen, _ := password.NewGenerator(input)
	username, err := gen.Generate(15, 0, 0, true, true)
	if err != nil {
		return nil, err
	}

	return ev_email.NewEmail(username, domain), nil
}

const (
	RandomRCPTStage = CloseStage + 1
	ConnectionStage = RandomRCPTStage + 1
)

type ClientGetter func(addr string) (*smtp.Client, error)

type CheckerInterface interface {
	Validate(mxs utils.MXs, email ev_email.EmailAddressInterface) SMTPError
}

func SimpleClientGetter(addr string) (*smtp.Client, error) {
	return smtp.Dial(addr)
}

type Checker struct {
	GetConn   ClientGetter
	Auth      smtp.Auth
	SendMail  SendMailInterface
	FromEmail ev_email.EmailAddressInterface
}

func (c Checker) Validate(mxs utils.MXs, email ev_email.EmailAddressInterface) SMTPError {
	var client *smtp.Client
	var err error
	var host string

	for _, mx := range mxs {
		host = fmt.Sprintf("%v:%v", mx.Host, DefaultSMTPPort)
		if client, err = c.GetConn(host); err == nil {
			break
		}
	}
	if client == nil {
		if err != nil {
			err = errors.New("smtp: connection was not created")
		}

		return NewSmtpError(ConnectionStage, err)
	}
	c.SendMail.SetClient(client)
	defer c.SendMail.Close()

	if err = c.SendMail.Hello(); err != nil {
		return NewSmtpError(HelloStage, err)
	}
	if err = c.SendMail.Auth(c.Auth); err != nil {
		return NewSmtpError(AuthStage, err)
	}

	err = c.SendMail.Mail(c.FromEmail.String())
	if err != nil {
		return NewSmtpError(MailStage, err)
	}

	if err = c.SendMail.RCPT(email.String()); err != nil {
		return NewSmtpError(RCPTStage, err)
	}

	rEmail, err := randomEmail(email.Domain())
	if err != nil {
		panic(err)
	}
	if err = c.SendMail.RCPT(rEmail.String()); err != nil {
		return NewSmtpError(RandomRCPTStage, err)
	}

	if err = c.SendMail.Quit(); err != nil {
		return NewSmtpError(QuitStage, err)
	}

	return nil
}
