package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	mockevmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"reflect"
	"testing"
)

func TestNewInput(t *testing.T) {
	type args struct {
		email     evmail.Address
		kvOptions []KVOption
	}

	kvOptions := []KVOption{
		{Name: OtherValidator, Option: 1},
		{Name: OtherValidator, Option: 3},
		{Name: SMTPValidatorName, Option: 2},
	}

	tests := []struct {
		name string
		args args
		want Input
	}{
		{
			name: "success",
			args: args{
				email:     mockevmail.GetValidTestEmail(),
				kvOptions: kvOptions,
			},
			want: &input{
				email: mockevmail.GetValidTestEmail(),
				options: map[ValidatorName]interface{}{
					OtherValidator:    3,
					SMTPValidatorName: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := make([]KVOption, 0)

			for _, opt := range tt.args.kvOptions {
				opts = append(opts, NewKVOption(opt.Name, opt.Option))
			}

			if got := NewInput(tt.args.email, opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
