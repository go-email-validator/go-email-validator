package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/disposable"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"reflect"
	"testing"
)

type mockDisposable struct {
	result bool
}

func (m mockDisposable) Disposable(_ ev_email.EmailAddressInterface) bool {
	return m.result
}

func getMockDisposable(result bool) mockDisposable {
	return mockDisposable{result}
}

func TestDisposableValidator_Validate(t *testing.T) {
	type fields struct {
		d disposable.Interface
	}
	type args struct {
		email ev_email.EmailAddressInterface
		in1   []ValidationResultInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ValidationResultInterface
	}{
		{
			name: "valid",
			fields: fields{
				d: getMockDisposable(false),
			},
			args: args{email: GetValidTestEmail()},
			want: NewValidatorResult(true, nil, nil, DisposableValidatorName),
		},
		{
			name: "invalid",
			fields: fields{
				d: getMockDisposable(true),
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
