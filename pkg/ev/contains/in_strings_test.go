package contains

import (
	"github.com/emirpasic/gods/sets/hashset"
	"reflect"
	"testing"
)

const (
	firstValue   = "first"
	secondValue  = "second"
	missingValue = "missing"
)

var twoStrings = []string{firstValue, secondValue}
var twoStringsInterface = []interface{}{firstValue, secondValue}
var setTwoStrings = NewSet(hashset.New(twoStringsInterface...))

func TestNewInStringsFromArray(t *testing.T) {
	type args struct {
		elements []string
	}

	tests := []struct {
		name string
		args args
		want InStrings
	}{
		{
			name: "",
			args: args{
				elements: twoStrings,
			},
			want: inStrings{setTwoStrings, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInStringsFromArray(tt.args.elements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInStringsFromArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inStrings_Contains(t *testing.T) {
	type fields struct {
		contains Interface
		maxLen   int
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "has " + firstValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   6,
			},
			args: args{
				value: firstValue,
			},
			want: true,
		},
		{
			name: "has " + secondValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   6,
			},
			args: args{
				value: secondValue,
			},
			want: true,
		},
		{
			name: "missing of " + missingValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   6,
			},
			args: args{
				value: missingValue,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := NewInStrings(tt.fields.contains, tt.fields.maxLen)
			if got := is.Contains(tt.args.value); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
