package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"reflect"
	"testing"
)

func TestDisposableValidator_Validate(t *testing.T) {
	type fields struct {
		d contains.InSet
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
			name: "valid",
			fields: fields{
				d: newMockContains(nil),
			},
			args: args{email: GetValidTestEmail()},
			want: NewValidatorResult(true, nil, nil, DisposableValidatorName),
		},
		{
			name: "invalid",
			fields: fields{
				d: newMockContains(validDomain),
			},
			args: args{email: GetValidTestEmail()},
			want: NewValidatorResult(false, []error{DisposableError{}}, nil, DisposableValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDisposableValidator(tt.fields.d)
			if got := d.Validate(tt.args.email, tt.args.in1...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
