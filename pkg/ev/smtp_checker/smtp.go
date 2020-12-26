package smtp_checker

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/log"
	"github.com/sethvargo/go-password/password"
	"net"
	"net/smtp"
)

const (
	DefaultEmail    = "user@example.org"
	DefaultSMTPPort = 25
)

func randomEmail(domain string) (ev_email.EmailAddress, error) {
	input := new(password.GeneratorInput)
	input.LowerLetters = password.LowerLetters + password.Digits

	gen, _ := password.NewGenerator(input)
	username, err := gen.Generate(15, 0, 0, true, true)
	if err != nil {
		return nil, err
	}

	return ev_email.NewEmailAddress(username, domain), nil
}

const (
	RandomRCPTStage = CloseStage + 1
	ConnectionStage = RandomRCPTStage + 1
)

// Direct DialFunc smtp.Dial
type DialFunc func(addr string) (*smtp.Client, error)

type Checker interface {
	Validate(mxs utils.MXs, email ev_email.EmailAddress) []error
}

type CheckerDTO struct {
	DialFunc  DialFunc
	SendMail  SendMail
	FromEmail ev_email.EmailAddress
}

func NewChecker(dto CheckerDTO) Checker {
	return checker{
		dialFunc:  dto.DialFunc,
		Auth:      nil,
		sendMail:  dto.SendMail,
		fromEmail: dto.FromEmail,
	}
}

/*
Some SMTP server send additional message and we should read it
2.1.5 for OK message
*/
type checker struct {
	dialFunc  DialFunc // use for get connection to smtp server
	Auth      smtp.Auth
	sendMail  SendMail
	fromEmail ev_email.EmailAddress
}

func (c checker) Validate(mxs utils.MXs, email ev_email.EmailAddress) (errs []error) {
	var client *smtp.Client
	var err error
	errs = make([]error, 0)
	var host string
	var e *net.OpError

	for _, mx := range mxs {
		host = fmt.Sprintf("%v:%v", mx.Host, DefaultSMTPPort)
		if client, err = c.dialFunc(host); err == nil {
			break
		}
		if !errors.As(err, &e) {
			log.Logger().Error(err)
		}
	}
	if client == nil {
		if err != nil {
			err = errors.New(fmt.Sprintf("smtp: connection was not created \n %s", err))
		}

		errs = append(errs, NewSmtpError(ConnectionStage, err))
		return
	}
	c.sendMail.SetClient(client)
	defer c.sendMail.Close()

	if err = c.sendMail.Hello(); err != nil {
		errs = append(errs, NewSmtpError(HelloStage, err))
		return
	}
	if err = c.sendMail.Auth(c.Auth); err != nil {
		errs = append(errs, NewSmtpError(AuthStage, err))
		return
	}

	err = c.sendMail.Mail(c.fromEmail.String())
	if err != nil {
		errs = append(errs, NewSmtpError(MailStage, err))
		return
	}

	rEmail, err := randomEmail(email.Domain())
	if err != nil {
		panic(err)
	}
	if err = c.sendMail.RCPT(rEmail.String()); err != nil {
		errs = append(errs, NewSmtpError(RandomRCPTStage, err))

		if err = c.sendMail.RCPT(email.String()); err != nil {
			errs = append(errs, NewSmtpError(RCPTStage, err))
		}
	}

	if err = c.sendMail.Quit(); err != nil {
		errs = append(errs, NewSmtpError(QuitStage, err))
	}

	return
}
