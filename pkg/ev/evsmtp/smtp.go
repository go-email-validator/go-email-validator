package evsmtp

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"github.com/go-email-validator/go-email-validator/pkg/log"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"net"
	"net/smtp"
)

// Configuration constants
const (
	ErrPrefix        = "evsmtp: "
	ErrConnectionMsg = ErrPrefix + "connection was not created \n %w"
	DefaultEmail     = "user@example.org"
	DefaultSMTPPort  = 25
	DefaultLocalName = "localhost"
)

// MXs is short alias for []*net.MX
type MXs = []*net.MX

// Constants of stages
const (
	RandomRCPTStage = CloseStage + 1
	ConnectionStage = RandomRCPTStage + 1
)

var (
	// ErrConnection is error of connection
	ErrConnection = NewError(ClientStage, errors.New(ErrConnectionMsg))
	// DefaultFromEmail is address, used as default From email
	DefaultFromEmail = evmail.FromString(DefaultEmail)
)

// DialFunc is function type to create smtpclient.SMTPClient
type DialFunc func(addr string) (smtpclient.SMTPClient, error)

// Dial generates smtpclient.SMTPClient (smtp.Client)
func Dial(addr string) (smtpclient.SMTPClient, error) {
	client, err := smtp.Dial(addr)
	return client, err
}

// RandomRCPTFunc is function for checking of Catching All
type RandomRCPTFunc func(email evmail.Address) (errs []error)

// RandomRCPT Need to realize of is-a relation (inheritance)
type RandomRCPT interface {
	Call(email evmail.Address) []error
	set(fn RandomRCPTFunc)
	get() RandomRCPTFunc
}

// ARandomRCPT is abstract realization of RandomRCPT
type ARandomRCPT struct {
	fn RandomRCPTFunc
}

// Call is calling of RandomRCPTFunc
func (a *ARandomRCPT) Call(email evmail.Address) []error {
	return a.fn(email)
}

func (a *ARandomRCPT) set(fn RandomRCPTFunc) {
	a.fn = fn
}

func (a *ARandomRCPT) get() RandomRCPTFunc {
	return a.fn
}

// Checker is SMTP validation
type Checker interface {
	Validate(mxs MXs, email evmail.Address) []error
}

// CheckerWithRandomRCPT is used for caching of RandomRCPT
type CheckerWithRandomRCPT interface {
	Checker
	RandomRCPT
}

// RandomEmail is function type to generate random email for checking of Catching All emails by RCPTs
type RandomEmail func(domain string) (evmail.Address, error)

func randomEmail(domain string) (evmail.Address, error) {
	gen, _ := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: password.LowerLetters + password.Digits,
	})
	username, err := gen.Generate(15, 0, 0, true, true)
	return evmail.NewEmailAddress(username, domain), err
}

// CheckerDTO is DTO for NewChecker
type CheckerDTO struct {
	DialFunc    DialFunc
	SendMail    SendMail
	FromEmail   evmail.Address
	LocalName   string
	RandomEmail RandomEmail
	Port        int
}

// NewChecker instantiates Checker
func NewChecker(dto CheckerDTO) Checker {
	if dto.DialFunc == nil {
		dto.DialFunc = Dial
	}

	if dto.SendMail == nil {
		dto.SendMail = NewSendMail(nil)
	}

	if dto.FromEmail == nil {
		dto.FromEmail = DefaultFromEmail
	}

	if dto.LocalName == "" {
		dto.LocalName = DefaultLocalName
	}

	if dto.RandomEmail == nil {
		dto.RandomEmail = randomEmail
	}

	if dto.Port == 0 {
		dto.Port = DefaultSMTPPort
	}

	c := checker{
		dialFunc:    dto.DialFunc,
		Auth:        nil,
		sendMail:    dto.SendMail,
		fromEmail:   dto.FromEmail,
		localName:   dto.LocalName,
		randomEmail: dto.RandomEmail,
		port:        dto.Port,
	}
	c.RandomRCPT = &ARandomRCPT{fn: c.randomRCPT}

	return c
}

/*
Some SMTP server send additional message and we should read it
2.1.5 for OK message
*/
type checker struct {
	RandomRCPT
	dialFunc    DialFunc // use for get connection to smtp server
	Auth        smtp.Auth
	sendMail    SendMail
	fromEmail   evmail.Address
	localName   string
	randomEmail RandomEmail
	port        int
}

func (c checker) Validate(mxs MXs, email evmail.Address) (errs []error) {
	var client interface{}
	var err error
	errs = make([]error, 0)
	var host string

	for _, mx := range mxs {
		host = fmt.Sprintf("%v:%v", mx.Host, c.port)
		if client, err = c.dialFunc(host); err == nil {
			break
		}
	}

	if err != nil {
		return append(errs, NewError(ConnectionStage, err))
	}

	if client == nil {
		return append(errs, ErrConnection)
	}

	c.sendMail.SetClient(client)
	needClose := true
	defer func() {
		if !needClose {
			return
		}
		if err = c.sendMail.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("sendMail.Close %v", err),
				zap.String("email", email.String()),
				zap.String("mxs", fmt.Sprint(mxs)),
			)
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

	if errsRandomRCPTs := c.RandomRCPT.Call(email); len(errsRandomRCPTs) > 0 {
		errs = append(errs, errsRandomRCPTs...)
		if errsRCPTs := c.sendMail.RCPTs([]string{email.String()}); len(errsRCPTs) > 0 {
			errs = append(errs, NewError(RCPTsStage, errsRCPTs[email.String()]))
		}
	}

	needClose = false
	if err = c.sendMail.Quit(); err != nil {
		errs = append(errs, NewError(QuitStage, err))
	}

	return
}

func (c checker) randomRCPT(email evmail.Address) (errs []error) {
	randomEmail, err := c.randomEmail(email.Domain())
	if err != nil {
		randomEmailErr := NewError(RandomRCPTStage, err)
		log.Logger().Error(fmt.Sprintf("generate random email: %v", randomEmailErr),
			zap.String("email", email.String()),
		)
		return append(errs, randomEmailErr)
	}

	if errsRCPTs := c.sendMail.RCPTs([]string{randomEmail.String()}); len(errsRCPTs) > 0 {
		return append(errs, NewError(RandomRCPTStage, errsRCPTs[randomEmail.String()]))
	}

	return
}

// RandomCacheKeyGetter is type of function to get cache key
type RandomCacheKeyGetter func(email evmail.Address) interface{}

// DefaultRandomCacheKeyGetter generates of cache key for RandomRCPT
func DefaultRandomCacheKeyGetter(email evmail.Address) interface{} {
	return email.Domain()
}

// NewCheckerCacheRandomRCPT creates Checker with caching of RandomRCPT calling
func NewCheckerCacheRandomRCPT(checker CheckerWithRandomRCPT, cache evcache.Interface, getKey RandomCacheKeyGetter) Checker {
	if getKey == nil {
		getKey = DefaultRandomCacheKeyGetter
	}

	c := &checkerCacheRandomRCPT{
		CheckerWithRandomRCPT: checker,
		randomRCPT:            &ARandomRCPT{fn: checker.get()},
		cache:                 cache,
		getKey:                getKey,
	}

	c.CheckerWithRandomRCPT.set(c.RandomRCPT)

	return c
}

type checkerCacheRandomRCPT struct {
	CheckerWithRandomRCPT
	randomRCPT RandomRCPT
	cache      evcache.Interface
	getKey     RandomCacheKeyGetter
}

func (c checkerCacheRandomRCPT) RandomRCPT(email evmail.Address) (errs []error) {
	key := c.getKey(email)
	resultInterface, err := c.cache.Get(key)
	if err == nil && resultInterface != nil {
		errs = *resultInterface.(*[]error)
	} else {
		errs = c.randomRCPT.Call(email)
		if err = c.cache.Set(key, ErrorsToEVSMTPErrors(errs)); err != nil {
			log.Logger().Error(fmt.Sprintf("cache RandomRCPT: %s", err),
				zap.String("email", email.String()),
				zap.String("key", fmt.Sprint(key)),
			)
		}
	}

	return errs
}
