package evsmtp

import (
	"net"
	"reflect"
	"testing"
)

func TestLookupMX(t *testing.T) {
	//evtests.FunctionalSkip(t)

	type args struct {
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    []*net.MX
		wantErr bool
	}{
		{
			name: "valid domain",
			args: args{
				domain: "gmail.com",
			},
			want: MXs{
				&net.MX{
					Host: "gmail-smtp-in.l.google.com.",
					Pref: 5,
				},
				&net.MX{
					Host: "alt1.gmail-smtp-in.l.google.com.",
					Pref: 10,
				},
				&net.MX{
					Host: "alt2.gmail-smtp-in.l.google.com.",
					Pref: 20,
				},
				&net.MX{
					Host: "alt3.gmail-smtp-in.l.google.com.",
					Pref: 30,
				},
				&net.MX{
					Host: "alt4.gmail-smtp-in.l.google.com.",
					Pref: 40,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid domain",
			args: args{
				domain: "someemailwithcannotexist.com",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LookupMX(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupMX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupMX() got = %v, want %v", got, tt.want)
			}
		})
	}
}
