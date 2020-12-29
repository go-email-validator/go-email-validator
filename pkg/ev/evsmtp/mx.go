package evsmtp

import (
	"net"
)

type FuncLookupMX func(domain string) ([]*net.MX, error)

func LookupMX(domain string) ([]*net.MX, error) {
	return net.LookupMX(domain)
}
