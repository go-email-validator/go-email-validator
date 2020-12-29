package contains

import (
	"github.com/emirpasic/gods/sets/hashset"
	"reflect"
	"testing"
)

const (
	firstValue   = "first"
	longValue    = "very_long_value_which_we_can_find_in_email"
	missingValue = "missing"
)

var twoStrings = []string{firstValue, longValue}
var twoStringsInterface = []interface{}{firstValue, longValue}
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
			name: "success",
			args: args{
				elements: twoStrings,
			},
			want: inStrings{setTwoStrings, len(longValue)},
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
		contains InSet
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
			name: "has equivalent of " + firstValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: firstValue,
			},
			want: true,
		},
		{
			name: "has " + firstValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: missingValue + firstValue + missingValue,
			},
			want: true,
		},
		{
			name: "has " + longValue + "in start",
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: longValue + missingValue,
			},
			want: true,
		},
		{
			name: "has " + longValue + "in end",
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: missingValue + longValue,
			},
			want: true,
		},
		{
			name: "missing of " + missingValue,
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: missingValue,
			},
			want: false,
		},
		{
			name: "empty value",
			fields: fields{
				contains: setTwoStrings,
				maxLen:   len(longValue),
			},
			args: args{
				value: "",
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
