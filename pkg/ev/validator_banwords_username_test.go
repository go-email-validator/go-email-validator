package ev

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/contains"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"reflect"
	"testing"
)

func Test_banWordsUsernameValidator_Validate(t *testing.T) {
	type fields struct {
		d contains.InStrings
	}
	type args struct {
		email evmail.Address
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ValidationResult
	}{
		{
			name: "username is valid",
			fields: fields{
				d: mockInString{
					t:    t,
					want: validEmail.Username(),
					ret:  false,
				},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(true, nil, nil, BanWordsUsernameValidatorName),
		},
		{
			name: "username is banned",
			fields: fields{
				d: mockInString{
					t:    t,
					want: validEmail.Username(),
					ret:  true,
				},
			},
			args: args{
				email: validEmail,
			},
			want: NewResult(false, utils.Errs(BanWordsUsernameError{}), nil, BanWordsUsernameValidatorName),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewBanWordsUsername(tt.fields.d)
			if got := w.Validate(NewInput(tt.args.email)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBanWordsUsernameError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: BanWordsUsernameErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ba := BanWordsUsernameError{}
			if got := ba.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
