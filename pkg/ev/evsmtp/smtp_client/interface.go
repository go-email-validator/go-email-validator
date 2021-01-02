package smtp_client

import (
	"crypto/tls"
	"io"
	"net/smtp"
)

type SMTPClient interface {
	Close() error
	Hello(localName string) error
	StartTLS(config *tls.Config) error
	TLSConnectionState() (state tls.ConnectionState, ok bool)
	Verify(addr string) error
	Auth(a smtp.Auth) error
	Mail(from string) error
	Rcpt(to string) error
	Data() (io.WriteCloser, error)
	Extension(ext string) (bool, string)
	Reset() error
	Noop() error
	Quit() error
}
