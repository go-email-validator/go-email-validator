package converter

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
)

func MX2String(MXs evsmtp.MXs) []string {
	var result = make([]string, len(MXs))
	for i, mx := range MXs {
		result[i] = mx.Host
	}

	return result
}
