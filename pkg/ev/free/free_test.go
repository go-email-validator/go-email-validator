package free

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"testing"
)

const freeDomain = "domain.com"

func TestSetFree_IsFree(t *testing.T) {
	type fields struct {
		set sets.Set
	}
	type args struct {
		email ev_email.EmailAddressInterface
	}

	set := hashset.New(freeDomain)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "is free",
			fields: fields{set},
			args:   args{ev_email.NewEmailAddress("", freeDomain)},
			want:   true,
		},
		{
			name:   "is not free",
			fields: fields{set},
			args:   args{ev_email.NewEmailAddress("", "")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSetFree(tt.fields.set)
			if got := s.IsFree(tt.args.email); got != tt.want {
				t.Errorf("IsFree() = %v, want %v", got, tt.want)
			}
		})
	}
}
