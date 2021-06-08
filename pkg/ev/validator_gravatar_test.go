package ev

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"net/url"
	"reflect"
	"testing"
)

const GravatarExistEmail = "beau@dentedreality.com.au"
const GravatarExistEmailURL = "https://www.gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50?d=404"

// TODO mocking Gravatar
func Test_gravatarValidator_Validate(t *testing.T) {
	evtests.FunctionalSkip(t)

	type args struct {
		email   evmail.Address
		options []KVOption
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
			want: NewGravatarValidationResult(
				GravatarExistEmailURL,
				NewValidResult(GravatarValidatorName).(*AValidationResult),
			),
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
		{
			name: "expired timeout",
			args: args{
				email:   evmail.FromString("some.none.exist@with.non.exist.domain"),
				results: []ValidationResult{NewValidResult(SyntaxValidatorName)},
				options: []KVOption{NewKVOption(
					GravatarValidatorName,
					NewGravatarOptions(GravatarOptionsDTO{Timeout: 1}),
				)},
			},
			want: gravatarGetError(&url.Error{
				Op:  "Head",
				URL: "https://www.gravatar.com/avatar/77996abfe12fc2141488a60b29aa4844?d=404",
				Err: errors.New("context deadline exceeded (Client.Timeout exceeded while awaiting headers)"),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errStr string
			var wantErrStr string

			w := NewGravatarValidator()
			gotInterface := w.Validate(NewInput(tt.args.email, tt.args.options...), tt.args.results...)

			got := gotInterface.(gravatarValidationResult)
			want := tt.want.(gravatarValidationResult)

			if len(got.errors) > 0 && len(want.errors) > 0 {
				if errOp, ok := got.errors[0].(*url.Error); ok && errOp.Err != nil {
					errStr = errOp.Err.Error()
					errOp.Err = nil
				}

				wantErrOp, ok := want.errors[0].(*url.Error)
				if ok && wantErrOp.Err != nil {
					wantErrStr = wantErrOp.Err.Error()
					wantErrOp.Err = nil
				}
			}

			if !reflect.DeepEqual(got, want) || got.URL() != want.URL() || errStr != wantErrStr {
				t.Errorf("Validate() = %v, want %v", gotInterface, tt.want)
			}
		})
	}
}

func Test_gravatarValidator_race_parallel(t *testing.T) {
	evtests.FunctionalSkip(t)

	w := NewGravatarValidator()
	wantStr := "context deadline exceeded (Client.Timeout exceeded while awaiting headers)"
	for i := 0; i < 100; i++ {
		email := evmail.NewEmailAddress(fmt.Sprintf("someNoneExistUserName%d", i), "someNonExists")
		t.Run(email.String(), func(t *testing.T) {
			t.Parallel()

			gotInterface := w.Validate(
				NewInput(email, NewKVOption(
					GravatarValidatorName,
					NewGravatarOptions(GravatarOptionsDTO{Timeout: 1}),
				)),
				NewValidResult(SyntaxValidatorName))

			got := gotInterface.(gravatarValidationResult)
			gotStr := got.errors[0].(*url.Error).Err.Error()

			if gotStr != wantStr {
				t.Errorf("Validate() = %v, wantStr %v", gotStr, wantStr)
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
