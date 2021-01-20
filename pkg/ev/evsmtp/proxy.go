package evsmtp

import (
	"context"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"github.com/go-email-validator/go-email-validator/pkg/proxifier"
	"github.com/tevino/abool"
	"h12.io/socks"
	"net"
	"net/smtp"
)

// DialFunc is function type to create smtpclient.SMTPClient
type DialFunc func(ctx context.Context, addr, proxyURL string) (smtpclient.SMTPClient, error)

var smtpNewClient = smtp.NewClient

var directDial = DirectDial

// DirectDial generates smtpclient.SMTPClient (smtp.Client)
func DirectDial(ctx context.Context, addr, proxyURL string) (smtpclient.SMTPClient, error) {
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, proxifier.TCPConnection, addr)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtpNewClient(conn, host)
}

var h12ioDial = socks.Dial

// H12IODial generates smtpclient.SMTPClient (smtp.Client) with proxy in socks.Dial
func H12IODial(ctx context.Context, addr, proxyURL string) (smtpclient.SMTPClient, error) {
	if proxyURL == "" {
		return directDial(ctx, addr, proxyURL)
	}
	var c net.Conn
	var client smtpclient.SMTPClient = nil
	var err error
	p := h12ioDial(proxyURL)

	done := make(chan struct{}, 1)
	needClose := abool.New()
	go func() {
		c, err = p("tcp", addr)
		defer func() {
			defer func() { close(done) }()

			if needClose.IsSet() && c != nil {
				c.Close()
			}
		}()

		if err != nil {
			return
		}

		host, _, _ := net.SplitHostPort(addr)
		client, err = smtpNewClient(c, host)
	}()

	select {
	case <-ctx.Done():
		needClose.Set()
		return nil, ctx.Err()
	case <-done:
		return client, err
	}
}
