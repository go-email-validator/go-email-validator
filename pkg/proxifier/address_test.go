package proxifier

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

var (
	addressFirst   = "addressFirst"
	addressSecond  = "addressSecond"
	addressThird   = "addressThird"
	addressInvalid = "@@@@%%...sdfd"
	simpleError    = errors.New("simpleError")
)

func getTestAddrsStr() []string {
	return []string{addressFirst, addressSecond}
}

func getAddrsTest(t *testing.T, addrsStr []string) []*Address {
	addrs, errs := getAddressesFromString(addrsStr)
	if len(errs) > 0 {
		t.Error(errs)
	}

	return addrs
}

func getAddrErrs(addrStr []string) (errs []error) {
	for _, addr := range addrStr {
		_, err := url.Parse(addr)
		if err != nil {
			errs = append(errs, fmt.Errorf(InvalidAddr, addr, err))
		}
	}

	return
}
