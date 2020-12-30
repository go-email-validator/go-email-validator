package proxifier

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	addressFirst         = "addressFirst"
	addressFirstWithPort = addressFirst + ":88"
	addressSecond        = "addressSecond"
	addressThird         = "addressThird"
	addressInvalid       = "@@@@%%...sdfd"
	simpleError          = errors.New("simpleError")
)

func getTestAddrsStr() []string {
	return []string{addressFirstWithPort, addressSecond}
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

type addressValue struct {
	key   interface{}
	value interface{}
}

func mapAddress(values ...addressValue) MapAddress {
	m := newMap()

	for _, value := range values {
		m.Put(value.key, value.value)
	}

	return m
}

func TestCreateCircleAddress(t *testing.T) {
	type args struct {
		i     int
		m     MapAddress
		addrs []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Address
	}{
		{
			name: "start from 0",
			args: args{
				i: 0,
				m: mapAddress(
					addressValue{
						key: "key1",
						value: &Address{
							url: "key1",
						},
					}, addressValue{
						key: "key2",
						value: &Address{
							url: "key2",
						},
					},
				),
				addrs: []interface{}{
					"key1", "key2",
				},
			},
			want: &Address{
				url: "key1",
			},
		},
		{
			name: "start from 1",
			args: args{
				i: 1,
				m: mapAddress(
					addressValue{
						key: "key1",
						value: &Address{
							url: "key1",
						},
					}, addressValue{
						key: "key2",
						value: &Address{
							url: "key2",
						},
					},
				),
				addrs: []interface{}{
					"key1", "key2",
				},
			},
			want: &Address{
				url: "key2",
			},
		},
		{
			name: "circle",
			args: args{
				i: 2,
				m: mapAddress(
					addressValue{
						key: "key1",
						value: &Address{
							url: "key1",
						},
					}, addressValue{
						key: "key2",
						value: &Address{
							url: "key2",
						},
					},
				),
				addrs: []interface{}{
					"key1", "key2",
				},
			},
			want: &Address{
				url: "key1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateCircleAddress(tt.args.i)(tt.args.m, tt.args.addrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCircleAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRandomAddress(t *testing.T) {
	type args struct {
		m     MapAddress
		addrs []interface{}
	}

	// Can failed because of slow computing
	randSeed = func(seed int64) {
		assert.Equal(t, time.Now().UnixNano()/10000, seed/10000)
	}

	randIntn = func(n int) int {
		return n - 1
	}

	tests := []struct {
		name string
		args args
		want *Address
	}{
		{
			name: "success",
			args: args{
				m: mapAddress(
					addressValue{
						key: "key1",
						value: &Address{
							url: "key1",
						},
					}, addressValue{
						key: "key2",
						value: &Address{
							url: "key2",
						},
					},
				),
				addrs: []interface{}{
					"key1", "key2",
				},
			},
			want: &Address{
				url: "key2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomAddress(tt.args.m, tt.args.addrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRandomAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
