package proxifier

import (
	"math/rand"
	"time"
)

// Address stores address for proxy List
type Address struct {
	url  string
	used uint
	ban  bool
}

// GetAddress is type of function, which is used like rotation strategy to return Address.
// There are default realizations: GetRandomAddress, GetFirstAddress, GetLastAddress and CreateCircleAddress.
type GetAddress func(MapAddress, []interface{}) *Address

var (
	randSeed = rand.Seed
	randIntn = rand.Intn
)

var timeUnixNano = time.Now().UnixNano

// GetRandomAddress is strategy, returning random Address
func GetRandomAddress(m MapAddress, addrs []interface{}) *Address {
	randSeed(timeUnixNano())
	addr, _ := m.Get(addrs[randIntn(len(addrs))])

	return addr.(*Address)
}

// GetFirstAddress is strategy, returning first Address
func GetFirstAddress(m MapAddress, addrs []interface{}) *Address {
	addr, _ := m.Get(addrs[0])

	return addr.(*Address)
}

// GetLastAddress is strategy, returning last Address
func GetLastAddress(m MapAddress, addrs []interface{}) *Address {
	addr, _ := m.Get(addrs[len(addrs)-1])

	return addr.(*Address)
}

// CreateCircleAddress is circle strategy for rotation
func CreateCircleAddress(i int) GetAddress {
	return func(m MapAddress, addrs []interface{}) *Address {
		if i >= len(addrs) {
			i = 0
		}
		i++

		addrKey := addrs[i-1]
		addr, _ := m.Get(addrKey)

		return addr.(*Address)
	}
}
