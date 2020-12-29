package evsmtp

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/log"
	"github.com/sethvargo/go-password/password"
	"net"
	"net/smtp"
)

const (
	ErrPrefix       = "evsmtp: "
	ErrConnection   = ErrPrefix + "connection was not created \n %w"
	DefaultEmail    = "user@example.org"
	DefaultSMTPPort = 25
)

type MXs = []*net.MX

const (
	RandomRCPTStage = CloseStage + 1
	ConnectionStage = RandomRCPTStage + 1
)

// Direct DialFunc smtp.Dial
type DialFunc func(addr string) (interface{}, error)

func Dial(addr string) (interface{}, error) {
	return smtp.Dial(addr)
}

type Checker interface {
	Validate(mxs MXs, email evmail.Address) []error
}

type RandomEmail func(domain string) (evmail.Address, error)

func randomEmail(domain string) (evmail.Address, error) {
	gen, _ := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: password.LowerLetters + password.Digits,
	})
	username, err := gen.Generate(15, 0, 0, true, true)
	return evmail.NewEmailAddress(username, domain), err
}

type CheckerDTO struct {
	DialFunc    DialFunc
	SendMail    SendMail
	FromEmail   evmail.Address
	LocalName   string
	RandomEmail RandomEmail
}

func NewChecker(dto CheckerDTO) Checker {
	if dto.LocalName == "" {
		dto.LocalName = "localhost"
	}

	if dto.RandomEmail == nil {
		dto.RandomEmail = randomEmail
	}

	return checker{
		dialFunc:    dto.DialFunc,
		Auth:        nil,
		sendMail:    dto.SendMail,
		fromEmail:   dto.FromEmail,
		localName:   dto.LocalName,
		randomEmail: dto.RandomEmail,
	}
}

/*
Some SMTP server send additional message and we should read it
2.1.5 for OK message
*/
type checker struct {
	dialFunc    DialFunc // use for get connection to smtp server
	Auth        smtp.Auth
	sendMail    SendMail
	fromEmail   evmail.Address
	localName   string
	randomEmail RandomEmail
}

// TODO improve logging, add fields and context
func (c checker) Validate(mxs MXs, email evmail.Address) (errs []error) {
	var client interface{}
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
			err = fmt.Errorf(ErrConnection, err)
		}

		return append(errs, NewError(ConnectionStage, err))
	}
	c.sendMail.SetClient(client)
	defer func() {
		err = c.sendMail.Close()
		if err != nil {
			log.Logger().Error(err)
		}
	}()

	if err = c.sendMail.Hello(c.localName); err != nil {
		errs = append(errs, NewError(HelloStage, err))
		return
	}
	if err = c.sendMail.Auth(c.Auth); err != nil {
		errs = append(errs, NewError(AuthStage, err))
		return
	}

	err = c.sendMail.Mail(c.fromEmail.String())
	if err != nil {
		errs = append(errs, NewError(MailStage, err))
		return
	}

	commonEmailRCPT := func() {
		if err = c.sendMail.RCPT(email.String()); err != nil {
			errs = append(errs, NewError(RCPTStage, err))
		}
	}
	rEmail, err := c.randomEmail(email.Domain())
	if err == nil {
		if err = c.sendMail.RCPT(rEmail.String()); err != nil {
			errs = append(errs, NewError(RandomRCPTStage, err))
			commonEmailRCPT()
		}
	} else {
		log.Logger().Error(NewError(RandomRCPTStage, err))
		commonEmailRCPT()
	}

	if err = c.sendMail.Quit(); err != nil {
		errs = append(errs, NewError(QuitStage, err))
	}

	return
}
