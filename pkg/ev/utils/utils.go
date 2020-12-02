package utils

import (
	"fmt"
	"go/types"
	"net"
	"reflect"
)

type MXs = []*net.MX

type Void struct{}

func GetVoid() Void {
	var member Void

	return member
}

type StringSet map[string]Void

func Async(f interface{}, args ...interface{}) <-chan []interface{} {
	ch := make(chan []interface{})
	go func() {
		defer close(ch)
		v := reflect.ValueOf(f)

		in := make([]reflect.Value, len(args))
		for k, param := range args {
			in[k] = reflect.ValueOf(param)
		}

		var values = v.Call(in)
		var result = make([]interface{}, len(values))

		for index, value := range values {
			result[index] = value.Interface()
		}

		ch <- result
	}()

	return ch
}

func GetInt(args []interface{}) (int, error) {
	if len(args) > 0 {
		arg := &args[0]

		if reflect.TypeOf(*arg).Kind() == reflect.Int {
			return (*arg).(int), nil
		}
		return 0, fmt.Errorf("\"%T\" is not int", *arg)
	}

	return 0, fmt.Errorf("empty []")
}

func GetString(i interface{}) (string, error) {
	switch v := i.(type) {
	case string:
		return v, nil

	case fmt.Stringer:
		return v.String(), nil
	}

	return "", fmt.Errorf("argument \"%s\" should be fmt.Stringer or string", i)
}

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
