package smtp_checker

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

// smtp_checker.SendMail
type SendMailInterface interface {
	SetClient(interface{})
	Client() interface{}
	Hello() error
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

type SendMail struct {
	client    *smtp.Client
	TLSConfig *tls.Config
}

func NewSendMail() SendMailInterface {
	return &SendMail{}
}

func (c *SendMail) SetClient(client interface{}) {
	c.client = client.(*smtp.Client)
}

func (c SendMail) Client() interface{} {
	return c.client
}

func (c SendMail) Hello() error {
	var err error
	if err = c.client.Hello("localhost"); err != nil && err.Error() != SMTPErrorHelloAfter {
		return err
	}
	return nil
}

func (c SendMail) Auth(a smtp.Auth) error {
	var err error

	if ok, _ := c.client.Extension("STARTTLS"); ok && c.TLSConfig != nil {
		if testHookStartTLS != nil {
			testHookStartTLS(c.TLSConfig)
		}
		if err = c.client.StartTLS(c.TLSConfig); err != nil {
			return err
		}
	}

	if a != nil {
		if ok, _ := c.client.Extension("AUTH"); !ok {
			return errors.New("smtp_checker: server doesn't support AUTH")
		}
		if err = c.client.Auth(a); err != nil {
			return err
		}
	}
	return nil
}

func (c SendMail) Mail(from string) error {
	var err error

	if err = c.client.Mail(from); err != nil {
		return err
	}
	return nil
}

func (c SendMail) RCPTs(addr []string) error {
	var err error

	for _, addr := range addr {
		if err = c.client.Rcpt(addr); err != nil {
			return err
		}
	}
	return nil
}

func (c SendMail) RCPT(addr string) error {
	var err error

	if err = c.client.Rcpt(addr); err != nil {
		return err
	}
	return nil
}

func (c SendMail) Data() (io.WriteCloser, error) {
	var err error

	w, err := c.client.Data()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (c SendMail) Write(w io.WriteCloser, msg []byte) error {
	var err error

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c SendMail) Quit() error {
	return c.client.Quit()
}

func (c SendMail) Close() error {
	return c.client.Close()
}
