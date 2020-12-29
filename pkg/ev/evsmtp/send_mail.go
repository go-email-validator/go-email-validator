package evsmtp

import (
	"crypto/tls"
	"errors"
	"io"
	"net/smtp"
)

type SendMailStage uint8

const (
	ClientStage SendMailStage = iota + 1
	HelloStage
	AuthStage
	MailStage
	RCPTsStage
	RCPTStage
	DataStage
	WriteStage
	QuitStage
	CloseStage
)

type SendMail interface {
	SetClient(interface{})
	Client() interface{}
	Hello(localName string) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	RCPTs(addr []string) error
	RCPT(addr string) error
	Data() (io.WriteCloser, error)
	Write(w io.WriteCloser, msg []byte) error
	Quit() error
	Close() error
}

var testHookStartTLS func(*tls.Config)

func NewSendMail() SendMail {
	return &sendMail{}
}

type sendMail struct {
	client    *smtp.Client
	TLSConfig *tls.Config
}

func (s *sendMail) SetClient(client interface{}) {
	s.client = client.(*smtp.Client)
}

func (s *sendMail) Client() interface{} {
	return s.client
}

func (s *sendMail) Hello(localName string) error {
	if err := s.client.Hello(localName); err != nil && err.Error() != ErrorHelloAfter {
		return err
	}
	return nil
}

func (s *sendMail) Auth(a smtp.Auth) error {
	if ok, _ := s.client.Extension("STARTTLS"); ok && s.TLSConfig != nil {
		if testHookStartTLS != nil {
			testHookStartTLS(s.TLSConfig)
		}
		if err := s.client.StartTLS(s.TLSConfig); err != nil {
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

func (s *sendMail) RCPTs(addr []string) error {
	for _, addr := range addr {
		if err := s.client.Rcpt(addr); err != nil {
			return err
		}
	}

	return nil
}

func (s *sendMail) RCPT(addr string) error {
	return s.client.Rcpt(addr)
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
