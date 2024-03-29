package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/free"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"reflect"
	"testing"
)

func Test_freeValidator_Validate(t *testing.T) {
	type fields struct {
		f contains.InSet
	}
	type args struct {
		email evmail.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ValidationResult
	}{
		{
			name: "email is not free",
			fields: fields{
				f: mockContains{t: t, want: validEmail.Domain(), ret: false},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(true, nil, nil, FreeValidatorName),
		},
		{
			name: "email is free",
			fields: fields{
				f: mockContains{t: t, want: validEmail.Domain(), ret: true},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(false, utils.Errs(FreeError{}), nil, FreeValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFreeValidator(tt.fields.f)
			if got := r.Validate(NewInput(tt.args.email)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFreeDefaultValidator(t *testing.T) {
	tests := []struct {
		name string
		want Validator
	}{
		{
			name: "success",
			want: NewFreeValidator(free.NewWillWhiteSetFree()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FreeDefaultValidator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FreeDefaultValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFreeError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: FreeErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := FreeError{}
			if got := fr.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
