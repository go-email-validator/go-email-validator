package proxifier

import (
	"context"
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"net"
	"net/smtp"
)

// Constants to choose type of connection.
const (
	TCPConnection = "tcp"
	UDPConnection = "udp"
)

// Dialer is interface
type Dialer interface {
	proxy.Dialer
	proxy.ContextDialer
}

// ProxyDialerFunc is type of function, which generates connection through out proxyURI
type ProxyDialerFunc func(proxyURI string) func(ctx context.Context, network string, addr string) (net.Conn, error)

// SocksDialContext returns Dial for socks with context
func SocksDialContext(proxyURI string) func(ctx context.Context, network string, addr string) (net.Conn, error) {
	p := socks.Dial(proxyURI)

	return func(ctx context.Context, network string, addr string) (c net.Conn, err error) {
		done := make(chan struct{}, 1)
		needClose := false
		go func() {
			c, err = p(network, addr)
			done <- struct{}{}
			defer func() {
				if needClose && c != nil {
					c.Close()
				}
			}()
		}()

		select {
		case <-ctx.Done():
			needClose = true
			return nil, ctx.Err()
		case <-done:
			return c, err
		}
	}
}

// NewProxyDialer returns proxy.Dialer based on List
func NewProxyDialer(list List, dialerFunc ProxyDialerFunc) Dialer {
	if dialerFunc == nil {
		dialerFunc = SocksDialContext
	}

	return &dialer{
		list:       list,
		dialerFunc: dialerFunc,
	}
}

type dialer struct {
	list       List
	dialerFunc ProxyDialerFunc
}

func (d *dialer) DialContext(ctx context.Context, network, addr string) (c net.Conn, err error) {
	var proxyAddr string
	err = errors.New("init")
	for err != nil {
		if proxyAddr != "" {
			d.list.Ban(proxyAddr)
		}
		var proxyErr error
		proxyAddr, proxyErr = d.list.GetAddress()
		if proxyErr != nil {
			return nil, proxyErr
		}

		c, err = d.dialerFunc(proxyAddr)(ctx, network, addr)
	}

	return c, nil
}

func (d *dialer) Dial(network, addr string) (c net.Conn, err error) {
	return d.DialContext(context.Background(), network, addr)
}

// SMTPDialler is a means to establish a connection for SMTP.
type SMTPDialler interface {
	DialContext(ctx context.Context, addr string) (smtpclient.SMTPClient, error)
	Dial(addr string) (smtpclient.SMTPClient, error)
}

// ProxySMTPDialer creates SMTP Dialer from addresses
func ProxySMTPDialer(addrs []string) (SMTPDialler, []error) {
	lst, err := NewListFromStrings(ListDTO{Addresses: addrs})
	return NewSMTPDialer(NewProxyDialer(lst, nil), ""), err
}

// NewSMTPDialer creates SMTP Dialer based on proxy.Dialer
func NewSMTPDialer(dialer Dialer, network string) SMTPDialler {
	if network == "" {
		network = TCPConnection
	}

	return &smtpDialer{
		dialer:  dialer,
		network: network,
	}
}

type smtpDialer struct {
	dialer  Dialer
	network string
}

var smtpNewClient = smtp.NewClient

func (p *smtpDialer) DialContext(ctx context.Context, addr string) (smtpclient.SMTPClient, error) {
	conn, err := p.dialer.DialContext(ctx, p.network, addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtpNewClient(conn, host)
}

func (p *smtpDialer) Dial(addr string) (smtpclient.SMTPClient, error) {
	return p.DialContext(context.Background(), addr)
}
