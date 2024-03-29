package ev

import (
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	mockevmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"github.com/stretchr/testify/require"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var (
	emptyEmail                     = evmail.NewEmailAddress("", "")
	emptyErrors                    = make([]error, 0)
	validEmail                     = evmail.FromString(mockevmail.ValidEmailString)
	invalidEmail                   = evmail.FromString("some%..@invalid.%.email")
	validMockValidator   Validator = mockValidator{result: true}
	inValidMockValidator Validator = mockValidator{result: false}
	errorSimple                    = errors.New("errorSimple")
	errorSimple2                   = errors.New("errorSimple2")
	validResult                    = NewResult(true, nil, nil, OtherValidator)
	invalidResult                  = NewResult(false, utils.Errs(newMockError()), nil, OtherValidator)
)

func sortErrors(errs []error) []error {
	sort.Slice(errs, func(l, r int) bool {
		return strings.Compare(errs[l].Error(), errs[r].Error()) >= 0
	})

	return errs
}

type mockContains struct {
	t    *testing.T
	want interface{}
	ret  bool
}

func (m mockContains) Contains(value interface{}) bool {
	require.Equal(m.t, value, m.want)

	return m.ret
}

type mockInString struct {
	t    *testing.T
	want interface{}
	ret  bool
}

func (m mockInString) Contains(value string) bool {
	require.Equal(m.t, value, m.want)

	return m.ret
}

func newMockError() error {
	return mockError{}
}

type mockError struct{}

func (mockError) Error() string {
	return "mockError"
}

const mockValidatorName ValidatorName = "mockValidatorName"

func newMockValidator(result bool) mockValidator {
	return mockValidator{result: result}
}

type mockValidator struct {
	result bool
	deps   []ValidatorName
}

func (m mockValidator) Validate(_ Input, _ ...ValidationResult) ValidationResult {
	var err error
	if !m.result {
		err = newMockError()
	}

	return NewResult(m.result, utils.Errs(err), nil, OtherValidator)
}

func (m mockValidator) GetDeps() []ValidatorName {
	return m.deps
}

type mockValidationResult struct {
	errs  []error
	warns []error
	name  ValidatorName
}

func (m mockValidationResult) IsValid() bool {
	return m.HasErrors()
}

func (m mockValidationResult) Errors() []error {
	return m.errs
}

func (m mockValidationResult) HasErrors() bool {
	return reflect.ValueOf(m.Errors()).Len() > 0
}

func (m mockValidationResult) Warnings() []error {
	return m.warns
}

func (m mockValidationResult) HasWarnings() bool {
	return reflect.ValueOf(m.Warnings()).Len() > 0
}

func (m mockValidationResult) ValidatorName() ValidatorName {
	return m.name
}

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestMockValidator(t *testing.T) {
	cases := []struct {
		validator mockValidator
		expected  ValidationResult
	}{
		{
			validator: newMockValidator(true),
			expected:  NewResult(true, nil, nil, OtherValidator),
		},
		{
			validator: newMockValidator(false),
			expected:  NewResult(false, utils.Errs(newMockError()), nil, OtherValidator),
		},
	}

	var emptyEmail evmail.Address
	for _, c := range cases {
		actual := c.validator.Validate(NewInput(emptyEmail))
		require.Equal(t, c.expected, actual)
	}
}

func TestAValidatorWithoutDeps(t *testing.T) {
	validator := AValidatorWithoutDeps{}

	require.Equal(t, emptyDeps, validator.GetDeps())
}

func TestNewValidatorResult(t *testing.T) {
	type args struct {
		isValid  bool
		errors   []error
		warnings []error
		name     ValidatorName
	}
	tests := []struct {
		name string
		args args
		want ValidationResult
	}{
		{
			name: "empty name",
			args: args{
				isValid:  true,
				errors:   nil,
				warnings: nil,
				name:     "",
			},
			want: &validationResult{true, nil, nil, OtherValidator},
		},
		{
			name: "invalid with errors and warnings",
			args: args{
				isValid:  false,
				errors:   []error{errorSimple},
				warnings: []error{errorSimple},
				name:     mockValidatorName,
			},
			want: &validationResult{false, []error{errorSimple}, []error{errorSimple}, mockValidatorName},
		},
		{
			name: "invalid",
			args: args{
				isValid:  false,
				errors:   nil,
				warnings: nil,
				name:     mockValidatorName,
			},
			want: &validationResult{false, nil, nil, mockValidatorName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResult(tt.args.isValid, tt.args.errors, tt.args.warnings, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorName_String(t *testing.T) {
	tests := []struct {
		name string
		v    ValidatorName
		want string
	}{
		{
			name: "success",
			v:    mockValidatorName,
			want: string(mockValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAValidationResult_HasWarnings(t *testing.T) {
	type fields struct {
		isValid  bool
		errors   []error
		warnings []error
		name     ValidatorName
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				warnings: utils.Errs(errorSimple),
			},
			want: true,
		},
		{
			name: "false empty",
			fields: fields{
				warnings: []error{},
			},
			want: false,
		},
		{
			name: "false nil",
			fields: fields{
				warnings: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AValidationResult{
				isValid:  tt.fields.isValid,
				errors:   tt.fields.errors,
				warnings: tt.fields.warnings,
				name:     tt.fields.name,
			}
			if got := a.HasWarnings(); got != tt.want {
				t.Errorf("HasWarnings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAValidationResult_ValidatorName(t *testing.T) {
	type fields struct {
		isValid  bool
		errors   []error
		warnings []error
		name     ValidatorName
	}
	tests := []struct {
		name   string
		fields fields
		want   ValidatorName
	}{
		{
			name: "success",
			fields: fields{
				name: OtherValidator,
			},
			want: OtherValidator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AValidationResult{
				isValid:  tt.fields.isValid,
				errors:   tt.fields.errors,
				warnings: tt.fields.warnings,
				name:     tt.fields.name,
			}
			if got := a.ValidatorName(); got != tt.want {
				t.Errorf("ValidatorName() = %v, want %v", got, tt.want)
			}
		})
	}
}
