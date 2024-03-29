package evmail

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"reflect"
	"testing"
)

const (
	defaultDomainInput   = "Domain"
	defaultDomain        = "domain"
	defaultUsernameInput = "Username"
	defaultUsername      = "username"
)

func defaultEmailString() string {
	return defaultUsername + AT + defaultDomain
}

func defaultEmailInputString() string {
	return defaultUsernameInput + AT + defaultDomainInput
}

func defaultEmail() Address {
	return NewEmailAddress(defaultUsername, defaultDomain)
}

const (
	emptyDomain   = ""
	emptyUsername = ""
)

func defaultFields() fields {
	return fields{username: defaultUsernameInput, domain: defaultDomainInput}
}

func emptyEmailString() string {
	return emptyUsername + AT + emptyDomain
}

func emptyEmail() Address {
	return NewEmailAddress(emptyUsername, emptyDomain)
}

type fields struct {
	username string
	domain   string
}

func emptyFields() fields { return fields{username: emptyUsername, domain: emptyDomain} }

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func TestEmailAddress_Domain(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "with domain",
			fields: defaultFields(),
			want:   defaultDomain,
		},
		{
			name:   "empty domain",
			fields: emptyFields(),
			want:   emptyDomain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmailAddress(tt.fields.username, tt.fields.domain)
			if got := e.Domain(); got != tt.want {
				t.Errorf("Domain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAddress_String(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "with email",
			fields: defaultFields(),
			want:   defaultEmailString(),
		},
		{
			name:   "empty email",
			fields: emptyFields(),
			want:   emptyEmailString(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmailAddress(tt.fields.username, tt.fields.domain)
			if got := e.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailAddress_Username(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "with username",
			fields: defaultFields(),
			want:   defaultUsername,
		},
		{
			name:   "empty username",
			fields: emptyFields(),
			want:   emptyUsername,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEmailAddress(tt.fields.username, tt.fields.domain)
			if got := e.Username(); got != tt.want {
				t.Errorf("Username() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailFromString(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want Address
	}{
		{
			name: "email",
			args: args{defaultEmailInputString()},
			want: defaultEmail(),
		},
		{
			name: "empty email",
			args: args{emptyEmailString()},
			want: emptyEmail(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromString(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeparatedEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "email",
			args:  args{defaultEmailString()},
			want:  defaultUsername,
			want1: defaultDomain,
		},
		{
			name:  "empty email",
			args:  args{emptyEmailString()},
			want:  emptyUsername,
			want1: emptyDomain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SeparateEmail(tt.args.email)
			if got != tt.want {
				t.Errorf("SeparateEmail() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SeparateEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestEmptyEmail(t *testing.T) {
	type args struct {
		val        Address
		defaultVal Address
	}

	emailTest := FromString("some@email.com")

	tests := []struct {
		name string
		args args
		want Address
	}{
		{
			name: "value",
			args: args{
				val:        defaultEmail(),
				defaultVal: emailTest,
			},
			want: defaultEmail(),
		},
		{
			name: "default",
			args: args{
				val:        nil,
				defaultVal: emailTest,
			},
			want: emailTest,
		},
		{
			name: "default nil",
			args: args{
				val:        nil,
				defaultVal: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyEmail(tt.args.val, tt.args.defaultVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmptyEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
