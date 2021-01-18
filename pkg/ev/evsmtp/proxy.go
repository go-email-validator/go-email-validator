package evsmtp

import (
	"context"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"github.com/go-email-validator/go-email-validator/pkg/proxifier"
	"h12.io/socks"
	"net"
	"net/smtp"
)

// DialFunc is function type to create smtpclient.SMTPClient
type DialFunc func(ctx context.Context, addr, proxyURL string) (smtpclient.SMTPClient, error)

// DirectDial generates smtpclient.SMTPClient (smtp.Client)
func DirectDial(ctx context.Context, addr, proxyURL string) (smtpclient.SMTPClient, error) {
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, proxifier.TCPConnection, addr)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// H12IODial generates smtpclient.SMTPClient (smtp.Client) with proxy in socks.Dial
func H12IODial(ctx context.Context, addr, proxyURL string) (client smtpclient.SMTPClient, err error) {
	var c net.Conn
	p := socks.Dial(proxyURL)

	done := make(chan struct{}, 1)
	needClose := false
	go func() {
		c, err = p("tcp", addr)
		defer func() {
			defer func() { done <- struct{}{} }()
			if needClose && c != nil {
				c.Close()
			}
		}()

		if err != nil {
			return
		}
		host, _, _ := net.SplitHostPort(addr)
		client, err = smtp.NewClient(c, host)
	}()

	select {
	case <-ctx.Done():
		needClose = true
		return nil, ctx.Err()
	case <-done:
		return
	}
}
