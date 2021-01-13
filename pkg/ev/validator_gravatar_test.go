package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"reflect"
	"testing"
)

const GravatarExistEmail = "beau@dentedreality.com.au"

// TODO mocking Gravatar
func Test_gravatarValidator_Validate(t *testing.T) {
	evtests.FunctionalSkip(t)

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
				results: []ValidationResult{NewValidResult(SyntaxValidatorName)},
			},
			want: NewValidResult(GravatarValidatorName),
		},
		{
			name: "invalid syntax",
			args: args{
				email:   evmail.FromString(""),
				results: []ValidationResult{syntaxGetError()},
			},
			want: gravatarGetError(NewDepsError()),
		},
		{
			name: "invalid in gravatar",
			args: args{
				email:   evmail.FromString("some.none.exist@with.non.exist.domain"),
				results: []ValidationResult{NewValidResult(SyntaxValidatorName)},
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

func Test_gravatarValidator_GetDeps(t *testing.T) {
	tests := []struct {
		name string
		want []ValidatorName
	}{
		{
			name: "success",
			want: []ValidatorName{SyntaxValidatorName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGravatarValidator()
			if got := g.GetDeps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDeps() = %v, want %v", got, tt.want)
			}
		})
	}
}
