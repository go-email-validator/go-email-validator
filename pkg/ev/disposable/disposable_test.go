package disposable

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"testing"
)

const disposableDomain = "disposable_domain"

func disposableDomainEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("", disposableDomain)
}

func notDisposableDomainEmail() ev_email.EmailAddressInterface {
	return ev_email.NewEmailAddress("", "")
}

func TestFuncDisposable_Disposable(t *testing.T) {
	type fields struct {
		funcChecker FuncChecker
	}
	type args struct {
		email ev_email.EmailAddressInterface
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "is disposable",
			fields: fields{func(_ ev_email.EmailAddressInterface) bool { return true }},
			args:   args{notDisposableDomainEmail()},
			want:   true,
		},
		{
			name:   "is not disposable",
			fields: fields{func(_ ev_email.EmailAddressInterface) bool { return false }},
			args:   args{notDisposableDomainEmail()},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFuncDisposable(tt.fields.funcChecker)
			if got := f.Disposable(tt.args.email); got != tt.want {
				t.Errorf("Disposable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDisposable_Disposable(t *testing.T) {
	type fields struct {
		set sets.Set
	}
	type args struct {
		email ev_email.EmailAddressInterface
	}

	set := hashset.New(disposableDomain)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "is disposable",
			fields: fields{set},
			args:   args{disposableDomainEmail()},
			want:   true,
		},
		{
			name:   "is not disposable",
			fields: fields{set},
			args:   args{notDisposableDomainEmail()},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSetDisposable(tt.fields.set)
			if got := s.Disposable(tt.args.email); got != tt.want {
				t.Errorf("Disposable() = %v, want %v", got, tt.want)
			}
		})
	}
}
