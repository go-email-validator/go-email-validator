package role

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"testing"
)

const defaultRole = "role"

func defaultRoleEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress(defaultRole, "")
}

func emptyRoleEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("", "")
}

func TestSetRole_HasRole(t *testing.T) {
	type fields struct {
		set sets.Set
	}
	type args struct {
		email ev_email.EmailAddressInterface
	}

	set := hashset.New(defaultRole)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "has role",
			fields: fields{set},
			args:   args{defaultRoleEmail()},
			want:   true,
		},
		{
			name:   "does not have role",
			fields: fields{set},
			args:   args{emptyRoleEmail()},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSetRole(tt.fields.set)
			if got := s.HasRole(tt.args.email); got != tt.want {
				t.Errorf("HasRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
