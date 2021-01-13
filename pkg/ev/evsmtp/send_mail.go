package evsmtp

import (
	"crypto/tls"
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"io"
	"net/smtp"
)

// SendMailStage is stage type of SendMail
type SendMailStage uint8

// Constants of stages
const (
	ClientStage SendMailStage = iota + 1
	HelloStage
	AuthStage
	MailStage
	RCPTsStage
	DataStage
	WriteStage
	QuitStage
	CloseStage
)

// SendMail is interface of custom realization as smtp.SendMail
type SendMail interface {
	SetClient(interface{})
	Client() interface{}
	Hello(localName string) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	RCPTs(addrs []string) map[string]error
	Data() (io.WriteCloser, error)
	Write(w io.WriteCloser, msg []byte) error
	Quit() error
	Close() error
}

var testHookStartTLS func(*tls.Config)

// NewSendMail instantiates SendMail
func NewSendMail(tlsConfig *tls.Config) SendMail {
	return &sendMail{
		tlsConfig: tlsConfig,
	}
}

type sendMail struct {
	client    smtpclient.SMTPClient
	tlsConfig *tls.Config
}

func (s *sendMail) SetClient(client interface{}) {
	s.client = client.(smtpclient.SMTPClient)
}

func (s *sendMail) Client() interface{} {
	return s.client
}

func (s *sendMail) Hello(localName string) error {
	return s.client.Hello(localName)
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
