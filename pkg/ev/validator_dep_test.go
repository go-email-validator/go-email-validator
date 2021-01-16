package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	mockevmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

type testSleep struct {
	sleep time.Duration
	mockValidator
	deps []ValidatorName
}

func (t testSleep) GetDeps() []ValidatorName {
	return t.deps
}

func (t testSleep) Validate(_ Interface, results ...ValidationResult) ValidationResult {
	time.Sleep(t.sleep)

	var isValid = true
	for _, result := range results {
		if !result.IsValid() {
			isValid = false
			break
		}
	}

	return NewDepValidatorResult(isValid && t.result, nil)
}

func TestDepValidator_Validate_Independent(t *testing.T) {
	email := mockevmail.GetValidTestEmail()
	strings := emptyDeps

	depValidator := NewDepValidator(
		map[ValidatorName]Validator{
			"test1": &testSleep{
				0,
				newMockValidator(true),
				strings,
			},
			"test2": &testSleep{
				0,
				newMockValidator(true),
				strings,
			},
			"test3": &testSleep{
				0,
				newMockValidator(false),
				strings,
			},
		},
	)

	v := depValidator.Validate(NewInput(email))
	require.False(t, v.IsValid())
}

func TestDepValidator_Validate_Dependent(t *testing.T) {
	email := mockevmail.GetValidTestEmail()
	strings := emptyDeps

	depValidator := NewDepValidator(map[ValidatorName]Validator{
		"test1": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			strings,
		},
		"test2": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			strings,
		},
		"test3": &testSleep{
			100 * time.Millisecond,
			newMockValidator(true),
			[]ValidatorName{"test1", "test2"},
		},
	},
	)

	v := depValidator.Validate(NewInput(email))
	require.True(t, v.IsValid())
}

func TestDepValidator_Validate_Full(t *testing.T) {
	evtests.FunctionalSkip(t)

	email := evmail.FromString(mockevmail.ValidEmailString)
	depValidator := NewDepBuilder(nil).Build()

	v := depValidator.Validate(NewInput(email))
	require.True(t, v.IsValid())
}

func Test_depValidationResult_Errors(t *testing.T) {
	type fields struct {
		isValid bool
		results DepResult
	}
	tests := []struct {
		name   string
		fields fields
		want   []error
	}{
		{
			name: "with Errors",
			fields: fields{
				isValid: false,
				results: DepResult{
					mockValidatorName:   mockValidationResult{errs: []error{errorSimple, errorSimple2}},
					SyntaxValidatorName: mockValidationResult{errs: []error{errorSimple2, errorSimple}},
				},
			},
			want: []error{errorSimple, errorSimple2, errorSimple2, errorSimple},
		},
		{
			name: "without Errors",
			fields: fields{
				isValid: false,
				results: DepResult{
					mockValidatorName:   mockValidationResult{},
					SyntaxValidatorName: mockValidationResult{},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDepValidatorResult(tt.fields.isValid, tt.fields.results)

			got := sortErrors(d.Errors())
			tt.want = sortErrors(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Errors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_depValidationResult_Warnings(t *testing.T) {
	type fields struct {
		isValid bool
		results DepResult
	}
	tests := []struct {
		name         string
		fields       fields
		wantWarnings []error
	}{
		{
			name: "with Warnings",
			fields: fields{
				isValid: false,
				results: DepResult{
					mockValidatorName:   mockValidationResult{warns: []error{errorSimple, errorSimple2}},
					SyntaxValidatorName: mockValidationResult{warns: []error{errorSimple2, errorSimple}},
				},
			},
			wantWarnings: []error{errorSimple, errorSimple2, errorSimple2, errorSimple},
		},
		{
			name: "without Warnings",
			fields: fields{
				isValid: false,
				results: DepResult{
					mockValidatorName:   mockValidationResult{},
					SyntaxValidatorName: mockValidationResult{},
				},
			},
			wantWarnings: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDepValidatorResult(tt.fields.isValid, tt.fields.results)

			gotWarnings := sortErrors(d.Warnings())
			tt.wantWarnings = sortErrors(tt.wantWarnings)

			if !reflect.DeepEqual(gotWarnings, tt.wantWarnings) {
				t.Errorf("Warnings() = %v, want %v", gotWarnings, tt.wantWarnings)
			}
		})
	}
}
