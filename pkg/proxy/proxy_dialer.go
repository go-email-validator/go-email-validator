package proxy

import (
	"errors"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"net"
	"net/smtp"
)

func NewProxyDialer(list ProxyList) proxy.Dialer {
	return &dialer{
		list: list,
	}
}

type dialer struct {
	list ProxyList
}

func (d *dialer) Dial(network, addr string) (c net.Conn, err error) {
	err = errors.New("init")
	for err != nil {
		proxyAddr, proxyErr := d.list.GetAddress()
		if proxyErr != nil {
			return nil, proxyErr
		}

		Dial := socks.Dial(proxyAddr)
		c, err = Dial(network, addr)
	}

	return c, nil
}

type SMTPDialler interface {
	Dial(addr string) (*smtp.Client, error)
}

func NewSMTPDialer(dialer proxy.Dialer, network string) SMTPDialler {
	if network == "" {
		network = "tcp"
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

func (p *smtpDialer) Dial(addr string) (*smtp.Client, error) {
	conn, err := p.dialer.Dial(p.network, addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
