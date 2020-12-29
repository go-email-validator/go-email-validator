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
			want: NewDepValidator(GetDefaultFactories()),
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
			d := NewDepBuilder(tt.fields.validators)
			if got := d.Build(); !reflect.DeepEqual(got, tt.want) {
				/*
					TODO find right way to compare struct with function.
					1. Use pointer for function
					2. Use InterfaceData()
				*/
				if tt.name != "nil" || len(got.(depValidator).deps) != len(tt.want.(depValidator).deps) {
					t.Errorf("Build() = %v, want %v", got, tt.want)
				}
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
		{
			name: "delete not exist element",
			fields: fields{
				validators: ValidatorMap{},
			},
			args: args{
				names: []ValidatorName{mockValidatorName, SyntaxValidatorName},
			},
			want: &DepBuilder{
				validators: ValidatorMap{},
			},
		},
		{
			name: "delete exist element",
			fields: fields{
				validators: ValidatorMap{
					mockValidatorName: newMockValidator(false),
					MXValidatorName:   newMockValidator(false)},
			},
			args: args{
				names: []ValidatorName{mockValidatorName, SyntaxValidatorName},
			},
			want: &DepBuilder{
				validators: ValidatorMap{MXValidatorName: newMockValidator(false)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDepBuilder(tt.fields.validators)
			if got := d.Delete(tt.args.names...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
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
		{
			name: "set",
			fields: fields{
				validators: ValidatorMap{},
			},
			args: args{
				name:      mockValidatorName,
				validator: newMockValidator(false),
			},
			want: &DepBuilder{
				validators: ValidatorMap{mockValidatorName: newMockValidator(false)},
			},
		},
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
