package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"reflect"
	"testing"
)

func Test_whiteListValidator_Validate(t *testing.T) {
	type fields struct {
		d contains.InSet
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
			name: "email is in white list",
			fields: fields{
				d: mockContains{t: t, want: validEmail.Domain(), ret: true},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(true, nil, nil, WhiteListDomainValidatorName),
		},
		{
			name: "email is not in white list",
			fields: fields{
				d: mockContains{t: t, want: validEmail.Domain(), ret: false},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(false, utils.Errs(WhiteListError{}), nil, WhiteListDomainValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWhiteListValidator(tt.fields.d)
			if got := w.Validate(NewInput(tt.args.email)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhiteListError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: WhiteListErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := WhiteListError{}
			if got := wh.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
