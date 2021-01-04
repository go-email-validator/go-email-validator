package proxifier

import (
	"os"
	"strings"
)

const EnvName = "PROXIES"

func EnvProxies() []string {
	str, exists := os.LookupEnv(EnvName)

	if !exists {
		return []string{}
	}

	addrs := strings.Split(str, " ")

	return addrs
}
