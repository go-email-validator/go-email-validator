package proxifier

import (
	"os"
	"strings"
)

// EnvName is environment name to store list of proxy addresses
const EnvName = "PROXIES"

// EnvProxies returns list of url for Address from .env file
func EnvProxies() []string {
	str, exists := os.LookupEnv(EnvName)

	if !exists {
		return []string{}
	}

	addrs := strings.Split(str, " ")

	return addrs
}
