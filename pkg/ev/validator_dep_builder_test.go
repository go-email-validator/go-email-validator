package ev

import (
	"reflect"
	"testing"
)

func TestDepBuilder_Build(t *testing.T) {
	type fields struct {
		validators ValidatorMap
	}
	tests := []struct {
		name   string
		fields fields
		want   Validator
	}{
		{
			name: "nil",
			fields: fields{
				validators: nil,
			},
			want: NewDepValidator(nil),
		},
		{
			name: "empty map",
			fields: fields{
				validators: ValidatorMap{},
			},
			want: NewDepValidator(ValidatorMap{}),
		},
		{
			name: "map",
			fields: fields{
				validators: ValidatorMap{
					mockValidatorName: newMockValidator(true),
				},
			},
			want: NewDepValidator(ValidatorMap{
				mockValidatorName: newMockValidator(true),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDepBuilder(&tt.fields.validators)
			if got := d.Build(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepBuilder_Delete(t *testing.T) {
	type fields struct {
		validators ValidatorMap
	}
	type args struct {
		names []ValidatorName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *DepBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DepBuilder{
				validators: tt.fields.validators,
			}
			if got := d.Delete(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepBuilder_Has(t *testing.T) {
	type fields struct {
		validators ValidatorMap
	}
	type args struct {
		names []ValidatorName
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DepBuilder{
				validators: tt.fields.validators,
			}
			if got := d.Has(tt.args.names...); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepBuilder_Set(t *testing.T) {
	type fields struct {
		validators ValidatorMap
	}
	type args struct {
		name      ValidatorName
		validator Validator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *DepBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DepBuilder{
				validators: tt.fields.validators,
			}
			if got := d.Set(tt.args.name, tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultFactories(t *testing.T) {
	tests := []struct {
		name string
		want *ValidatorMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultFactories(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefaultFactories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDepBuilder(t *testing.T) {
	type args struct {
		validators *ValidatorMap
	}
	tests := []struct {
		name string
		args args
		want *DepBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDepBuilder(tt.args.validators); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDepBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
