package evsmtp

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"io"
	"net/smtp"
	"sync"
)

// SendMailStage is stage type of SendMail
type SendMailStage uint8

// SafeSendMailStage is thread safe SendMailStage
type SafeSendMailStage struct {
	SendMailStage
	mu sync.RWMutex
}

// Set sets stage
func (s *SafeSendMailStage) Set(val SendMailStage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SendMailStage = val
}

// Get returns stage
func (s *SafeSendMailStage) Get() SendMailStage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SendMailStage
}

// Constants of stages
const (
	ClientStage SendMailStage = iota + 1
	HelloStage
	AuthStage
	MailStage
	RCPTsStage
	QuitStage
	CloseStage
)

// SendMail is interface of custom realization as smtp.SendMail
type SendMail interface {
	Client() smtpclient.SMTPClient
	Hello(helloName string) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	RCPTs(addrs []string) map[string]error
	Data() (io.WriteCloser, error)
	Write(w io.WriteCloser, msg []byte) error
	Quit() error
	Close() error
}

var testHookStartTLS func(*tls.Config)

// SendMailDialerFactory is factory for SendMail with dialing
type SendMailDialerFactory func(ctx context.Context, host string, opts Options) (SendMail, error)

// NewSendMailFactory creates SendMailDialerFactory
func NewSendMailFactory(dialFunc DialFunc, tlsConfig *tls.Config) SendMailDialerFactory {
	return NewSendMailCustom(dialFunc, tlsConfig, NewSendMail)
}

// SendMailFactory is factory for SendMail
type SendMailFactory func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail

// NewSendMailCustom creates SendMailFactory with dialing and customization calling of SendMailFactory
func NewSendMailCustom(dialFunc DialFunc, tlsConfig *tls.Config, factory SendMailFactory) SendMailDialerFactory {
	return func(ctx context.Context, host string, opts Options) (SendMail, error) {
		conn, err := dialFunc(ctx, host, opts.Proxy())
		if err != nil {
			return nil, err
		}

		return factory(conn, tlsConfig), nil
	}
}

// NewSendMail instantiates SendMail
func NewSendMail(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
	return &sendMail{
		client:    client,
		tlsConfig: tlsConfig,
	}
}

type sendMail struct {
	client    smtpclient.SMTPClient
	tlsConfig *tls.Config
}

func (s *sendMail) Client() smtpclient.SMTPClient {
	return s.client
}

func (s *sendMail) Hello(helloName string) error {
	return s.client.Hello(helloName)
}

func (s *sendMail) Auth(a smtp.Auth) error {
	if ok, _ := s.client.Extension("STARTTLS"); ok && s.tlsConfig != nil {
		if testHookStartTLS != nil {
			testHookStartTLS(s.tlsConfig)
		}
		if err := s.client.StartTLS(s.tlsConfig); err != nil {
			return err
		}
	}

	if a != nil {
		if ok, _ := s.client.Extension("AUTH"); !ok {
			return errors.New("smtp_checker: server doesn't support AUTH")
		}
		if err := s.client.Auth(a); err != nil {
			return err
		}
	}
	return nil
}

func (s *sendMail) Mail(from string) error {
	return s.client.Mail(from)
}

func (s *sendMail) RCPTs(addrs []string) map[string]error {
	errs := make(map[string]error)

	for _, addr := range addrs {
		if err := s.client.Rcpt(addr); err != nil {
			errs[addr] = err
		}
	}

	return errs
}

func (s *sendMail) Data() (io.WriteCloser, error) {
	return s.client.Data()
}

func (s *sendMail) Write(w io.WriteCloser, msg []byte) error {
	if _, err := w.Write(msg); err != nil {
		return err
	}

	return w.Close()
}

func (s *sendMail) Quit() error {
	return s.client.Quit()
}

func (s *sendMail) Close() error {
	return s.client.Close()
}
