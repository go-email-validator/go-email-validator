package domain

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"testing"
)

const defaultDomain = "default_domain"

func defaultDomainEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("", defaultDomain)
}

func emptyDomainEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("", "")
}

func Test_setDomain_Contains(t *testing.T) {
	type fields struct {
		set sets.Set
	}
	type args struct {
		email ev_email.EmailAddressInterface
	}

	set := hashset.New(defaultDomain)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "has domain",
			fields: fields{set},
			args:   args{defaultDomainEmail()},
			want:   true,
		},
		{
			name:   "does not have domain",
			fields: fields{set},
			args:   args{emptyDomainEmail()},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSetDomain(tt.fields.set)
			if got := s.Contains(tt.args.email); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
