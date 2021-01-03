package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"reflect"
	"regexp"
	"testing"
)

func Test_syntaxValidator_Validate(t *testing.T) {
	type fields struct{}
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
			name:   "success",
			fields: fields{},
			args: args{
				email: validEmail,
			},
			want: NewValidResult(SyntaxValidatorName),
		},
		{
			name:   "invalid",
			fields: fields{},
			args: args{
				email: invalidEmail,
			},
			want: syntaxGetError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sy := NewSyntaxValidator()
			if got := sy.Validate(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syntaxRegexValidator_Validate(t *testing.T) {
	invalidRegExp := regexp.MustCompile("^@$")

	type fields struct {
		AValidatorWithoutDeps AValidatorWithoutDeps
		emailRegex            *regexp.Regexp
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
			name: "success with default regex",
			fields: fields{
				emailRegex: nil,
			},
			args: args{
				email: validEmail,
			},
			want: NewValidResult(SyntaxValidatorName),
		},
		{
			name: "invalid with default regex",
			fields: fields{
				emailRegex: defaultEmailRegex,
			},
			args: args{
				email: invalidEmail,
			},
			want: syntaxGetError(),
		},
		{
			name: "invalid with custom regex",
			fields: fields{
				emailRegex: invalidRegExp,
			},
			args: args{
				email: validEmail,
			},
			want: syntaxGetError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSyntaxRegexValidator(tt.fields.emailRegex)
			if got := s.Validate(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
