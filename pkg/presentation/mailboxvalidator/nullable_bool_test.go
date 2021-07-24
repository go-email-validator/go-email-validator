package mailboxvalidator

import (
	"reflect"
	"testing"
)

func TestEmptyBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		emptyBool EmptyBool
		want      []byte
		wantErr   bool
	}{
		{name: "true", emptyBool: NewEmptyBool(true), want: []byte("true"), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.emptyBool
			got, err := e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyBool_ToString(t *testing.T) {
	tests := []struct {
		name      string
		emptyBool EmptyBool
		want      string
	}{
		{name: "true", emptyBool: NewEmptyBool(true), want: "True"},
		{name: "false", emptyBool: NewEmptyBool(false), want: "False"},
		{name: "nil", emptyBool: NewEmptyBoolWithNil(), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EmptyBool{
				bool: tt.emptyBool.bool,
			}
			if got := e.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantResult EmptyBool
	}{
		{name: "true", args: args{value: "True"}, wantResult: NewEmptyBool(true)},
		{name: "false", args: args{value: "False"}, wantResult: NewEmptyBool(false)},
		{name: "nil", args: args{value: ""}, wantResult: NewEmptyBoolWithNil()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := ToBool(tt.args.value); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ToBool() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
