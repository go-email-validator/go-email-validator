package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"reflect"
	"testing"
)

func Test_roleValidator_Validate(t *testing.T) {
	type fields struct {
		r contains.InSet
	}
	type args struct {
		email evmail.Address
		in1   []ValidationResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ValidationResult
	}{
		{
			name: "email has not role",
			fields: fields{
				r: mockContains{t: t, want: validEmail.Username(), ret: false},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(true, nil, nil, RoleValidatorName),
		},
		{
			name: "email has role",
			fields: fields{
				r: mockContains{t: t, want: validEmail.Username(), ret: true},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(false, utils.Errs(RoleError{}), nil, RoleValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRoleValidator(tt.fields.r)
			if got := r.Validate(NewInput(tt.args.email), tt.args.in1...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
