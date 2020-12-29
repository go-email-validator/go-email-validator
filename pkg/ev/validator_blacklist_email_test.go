package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"reflect"
	"testing"
)

func Test_blackListEmailsValidator_Validate(t *testing.T) {
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
			name: "email is valid",
			fields: fields{
				d: mockContains{t: t, want: validEmail.String(), ret: false},
			},
			args: args{
				email: validEmail,
			},
			want: NewValidatorResult(true, nil, nil, BlackListEmailsValidatorName),
		},
		{
			name: "email is invalid",
			fields: fields{
				d: mockContains{t: t, want: validEmail.String(), ret: true},
			},
			args: args{
				email: validEmail,
			},
			want: NewValidatorResult(false, utils.Errs(BlackListEmailsError{}), nil, BlackListEmailsValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewBlackListEmailsValidator(tt.fields.d)
			if got := w.Validate(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
