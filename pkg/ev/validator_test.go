package ev

import (
	"errors"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

var (
	emptyEmail                     = evmail.NewEmailAddress("", "")
	emptyErrors                    = make([]error, 0)
	validEmail                     = evmail.FromString(validEmailString)
	validMockValidator   Validator = mockValidator{result: true}
	inValidMockValidator Validator = mockValidator{result: false}
	simpleError                    = errors.New("simpleError")
	simpleError2                   = errors.New("simpleError2")
	sortErrors                     = func(errs []error) func(l, r int) bool {
		return func(l, r int) bool {
			return strings.Compare(errs[l].Error(), errs[l].Error()) > 0
		}
	}
)

type mockContains struct {
	t    *testing.T
	want interface{}
	ret  bool
}

func (m mockContains) Contains(value interface{}) bool {
	assert.Equal(m.t, value, m.want)

	return m.ret
}

type mockInString struct {
	t    *testing.T
	want interface{}
	ret  bool
}

func (m mockInString) Contains(value string) bool {
	assert.Equal(m.t, value, m.want)

	return m.ret
}

func newMockError() error {
	return mockError{}
}

type mockError struct {
	utils.Err
}

const mockValidatorName ValidatorName = "mockValidatorName"

func newMockValidator(result bool) mockValidator {
	return mockValidator{result: result}
}

type mockValidator struct {
	result bool
	AValidatorWithoutDeps
}

func (m mockValidator) Validate(_ evmail.Address, _ ...ValidationResult) ValidationResult {
	var err error
	if !m.result {
		err = newMockError()
	}

	return NewValidatorResult(m.result, utils.Errs(err), nil, OtherValidator)
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
			expected:  NewValidatorResult(true, nil, nil, OtherValidator),
		},
		{
			validator: newMockValidator(false),
			expected:  NewValidatorResult(false, utils.Errs(newMockError()), nil, OtherValidator),
		},
	}

	var emptyEmail evmail.Address
	for _, c := range cases {
		actual := c.validator.Validate(emptyEmail)
		assert.Equal(t, c.expected, actual)
	}
}

func TestAValidatorWithoutDeps(t *testing.T) {
	validator := AValidatorWithoutDeps{}

	assert.Equal(t, emptyDeps, validator.GetDeps())
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
				errors:   []error{simpleError},
				warnings: []error{simpleError},
				name:     mockValidatorName,
			},
			want: &validationResult{false, []error{simpleError}, []error{simpleError}, mockValidatorName},
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
			if got := NewValidatorResult(tt.args.isValid, tt.args.errors, tt.args.warnings, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidatorResult() = %v, want %v", got, tt.want)
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
