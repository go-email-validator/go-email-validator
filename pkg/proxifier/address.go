package proxifier

import (
	"math/rand"
	"time"
)

type Address struct {
	url  string
	used uint
	ban  bool
}

// interface{} consist of string.
type GetAddress func(MapAddress, []interface{}) *Address

var (
	randSeed = rand.Seed
	randIntn = rand.Intn
)

func GetRandomAddress(m MapAddress, addrs []interface{}) *Address {
	randSeed(time.Now().UnixNano())
	addr, _ := m.Get(addrs[randIntn(len(addrs))])

	return addr.(*Address)
}

func GetFirstAddress(m MapAddress, addrs []interface{}) *Address {
	addr, _ := m.Get(addrs[0])

	return addr.(*Address)
}

func GetLastAddress(m MapAddress, addrs []interface{}) *Address {
	addr, _ := m.Get(addrs[len(addrs)-1])

	return addr.(*Address)
}

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
