package proxy_list

import (
	"math/rand"
	"time"
)

type Address struct {
	url       string
	used      uint
	lastUsing time.Duration
	ban       bool
}

// []interface{} - []string
type GetAddress func(MapAddress, []interface{}) *Address

func GetRandomAddress(m MapAddress, addrs []interface{}) *Address {
	rand.Seed(time.Now().UnixNano())
	addr, _ := m.Get(addrs[rand.Intn(len(addrs))])
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
		addr, _ := m.Get(addrs[len(addrs)-1])
		return addr.(*Address)
	}
}
