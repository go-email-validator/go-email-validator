package evsmtp

import (
	"net"
)

// FuncLookupMX returns MXs
type FuncLookupMX func(domain string) (MXs, error)

// LookupMX is default realization for looking net.MX
func LookupMX(domain string) (MXs, error) {
	return net.LookupMX(domain)
}
