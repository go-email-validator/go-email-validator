package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	mock_evmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"github.com/stretchr/testify/require"
	"net"
	"reflect"
	"testing"
)

func mockLookupMX(t *testing.T, domainExpected string, ret evsmtp.MXs, err error) evsmtp.FuncLookupMX {
	return func(domain string) ([]*net.MX, error) {
		require.Equal(t, domainExpected, domain)

		return ret, err
	}
}

func Test_mxValidator_Validate(t *testing.T) {
	type fields struct {
		lookupMX evsmtp.FuncLookupMX
	}
	type args struct {
		email evmail.Address
	}

	mxs := evsmtp.MXs{&net.MX{}}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   ValidationResult
	}{
		{
			name: "existed domain",
			fields: fields{
				lookupMX: mockLookupMX(t, validEmail.Domain(), mxs, nil),
			},
			args: args{
				email: validEmail,
			},
			want: NewMXValidationResult(
				mxs,
				NewResult(true, nil, nil, MXValidatorName).(*AValidationResult),
			),
		},
		{
			name: "empty mx list",
			fields: fields{
				lookupMX: mockLookupMX(t, validEmail.Domain(), nil, nil),
			},
			args: args{
				email: validEmail,
			},
			want: NewMXValidationResult(
				nil,
				NewResult(false, utils.Errs(EmptyMXsError{}), nil, MXValidatorName).(*AValidationResult),
			),
		},
		{
			name: "unexisted domain",
			fields: fields{
				lookupMX: mockLookupMX(t, validEmail.Domain(), nil, simpleError),
			},
			args: args{
				email: validEmail,
			},
			want: NewMXValidationResult(
				nil,
				NewResult(false, utils.Errs(simpleError), nil, MXValidatorName).(*AValidationResult),
			),
		},
		{
			name: "unexisted domain with mxs",
			fields: fields{
				lookupMX: mockLookupMX(t, validEmail.Domain(), mxs, simpleError),
			},
			args: args{
				email: validEmail,
			},
			want: NewMXValidationResult(
				mxs,
				NewResult(false, utils.Errs(simpleError), nil, MXValidatorName).(*AValidationResult),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewMXValidator(tt.fields.lookupMX)
			if got := v.Validate(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSMTPValidator_Validate_MX(b *testing.B) {
	email := evmail.FromString(mock_evmail.ValidEmailString)

	depValidator := NewDepValidator(
		map[ValidatorName]Validator{
			SyntaxValidatorName: DefaultNewMXValidator(),
			MXValidatorName:     NewSyntaxValidator(),
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		depValidator.Validate(email)
	}
}
