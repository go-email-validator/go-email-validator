package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"reflect"
	"testing"
)

// TODO mocking Gravatar
func Test_gravatarValidator_Validate(t *testing.T) {
	type args struct {
		email   ev_email.EmailAddress
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
				email:   ev_email.EmailFromString("beau@dentedreality.com.au"),
				results: []ValidationResult{NewValidValidatorResult(SyntaxValidatorName)},
			},
			want: NewValidValidatorResult(GravatarValidatorName),
		},
		{
			name: "invalid syntax",
			args: args{
				email:   ev_email.EmailFromString(""),
				results: []ValidationResult{syntaxGetError()},
			},
			want: gravatarGetError(),
		},
		{
			name: "invalid in gravatar",
			args: args{
				email:   ev_email.EmailFromString("some.none.exist@with.non.exist.domain"),
				results: []ValidationResult{NewValidValidatorResult(SyntaxValidatorName)},
			},
			want: gravatarGetError(),
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
