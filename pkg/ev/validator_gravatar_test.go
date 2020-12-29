package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"reflect"
	"testing"
)

const GravatarExistEmail = "beau@dentedreality.com.au"

// TODO mocking Gravatar
func Test_gravatarValidator_Validate(t *testing.T) {
	type args struct {
		email   evmail.Address
		results []ValidationResult
	}

	tests := []struct {
		name string
		args args
		want ValidationResult
	}{
		{
			name: "valid",
			args: args{
				email:   evmail.FromString(GravatarExistEmail),
				results: []ValidationResult{NewValidValidatorResult(SyntaxValidatorName)},
			},
			want: NewValidValidatorResult(GravatarValidatorName),
		},
		{
			name: "invalid syntax",
			args: args{
				email:   evmail.FromString(""),
				results: []ValidationResult{syntaxGetError()},
			},
			want: gravatarGetError(DepsError{}),
		},
		{
			name: "invalid in gravatar",
			args: args{
				email:   evmail.FromString("some.none.exist@with.non.exist.domain"),
				results: []ValidationResult{NewValidValidatorResult(SyntaxValidatorName)},
			},
			want: gravatarGetError(GravatarError{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewGravatarValidator()
			if got := w.Validate(tt.args.email, tt.args.results...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
