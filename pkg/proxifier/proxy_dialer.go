package proxifier

import (
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtp_client"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"net"
	"net/smtp"
)

const (
	TCPConnection = "tcp"
	UDPConnection = "udp"
)

func NewProxyDialer(list List, dialerFunc ProxyDialerFunc) proxy.Dialer {
	if dialerFunc == nil {
		dialerFunc = socks.Dial
	}

	return &dialer{
		list:       list,
		dialerFunc: dialerFunc,
	}
}

type ProxyDialerFunc func(proxyURI string) func(network string, addr string) (net.Conn, error)

type dialer struct {
	list       List
	dialerFunc ProxyDialerFunc
}

func (d *dialer) Dial(network, addr string) (c net.Conn, err error) {
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

		c, err = d.dialerFunc(proxyAddr)(network, addr)
	}

	return c, nil
}

type SMTPDialler interface {
	Dial(addr string) (smtp_client.SMTPClient, error)
}

func ProxySmtpDialer(addrs []string) (SMTPDialler, []error) {
	lst, err := NewListFromStrings(ListDTO{Addresses: addrs})
	return NewSMTPDialer(NewProxyDialer(lst, nil), ""), err
}

func NewSMTPDialer(dialer proxy.Dialer, network string) SMTPDialler {
	if network == "" {
		network = TCPConnection
	}

	return &smtpDialer{
		dialer:  dialer,
		network: network,
	}
}

type smtpDialer struct {
	dialer  proxy.Dialer
	network string
}

var smtpNewClient = smtp.NewClient

func (p *smtpDialer) Dial(addr string) (smtp_client.SMTPClient, error) {
	conn, err := p.dialer.Dial(p.network, addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtpNewClient(conn, host)
}
