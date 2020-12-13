package utils

import (
	"fmt"
	"go/types"
	"net"
)

type MXs = []*net.MX

func RangeLen(i interface{}) int {
	if i == nil {
		return 0
	}

	switch i.(type) {
	case types.Array:
	case types.Slice:
		return len(i.([]interface{}))
	case types.Map:
		return len(i.(map[interface{}]interface{}))
	}
	panic(fmt.Sprintf("interface should be \"array\", \"slice\" or \"map\""))
}

func abstractFunc() interface{} {
	panic("implement me")
}
