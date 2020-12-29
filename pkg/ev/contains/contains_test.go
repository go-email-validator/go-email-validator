package contains

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"testing"
)

const defaultValue = "defaultValue"

func getFuncChecker(valueToCheck interface{}) FuncChecker {
	return func(val interface{}) bool { return val == valueToCheck }
}

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func Test_funcContains_Contains(t *testing.T) {
	type fields struct {
		funcChecker FuncChecker
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "contains",
			fields: fields{getFuncChecker(defaultValue)},
			args:   args{defaultValue},
			want:   true,
		},
		{
			name:   "is not free",
			fields: fields{getFuncChecker(defaultValue)},
			args:   args{nil},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFunc(tt.fields.funcChecker)
			if got := f.Contains(tt.args.value); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setContains_Contains(t *testing.T) {
	type fields struct {
		set sets.Set
	}
	type args struct {
		value interface{}
	}

	set := hashset.New(defaultValue)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "contains",
			fields: fields{set},
			args:   args{defaultValue},
			want:   true,
		},
		{
			name:   "is not free",
			fields: fields{set},
			args:   args{nil},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSet(tt.fields.set)
			if got := s.Contains(tt.args.value); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
