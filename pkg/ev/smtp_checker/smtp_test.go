package smtp_checker

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/test_utils"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/proxy_list"
	"github.com/stretchr/testify/assert"
	"net"
	"net/smtp"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	test_utils.TestMain(m)
}

func TestChecker_Validate(t *testing.T) {
	test_utils.FunctionalSkip(t)

	type fields struct {
		GetConn   DialFunc
		Auth      smtp.Auth
		SendMail  SendMail
		FromEmail ev_email.EmailAddress
	}
	type args struct {
		mxs   utils.MXs
		email ev_email.EmailAddress
	}

	emailString := "y-numata@senko.ed.jp"
	emailString = "asd@tradepro.net"

	proxyList, _ := proxy_list.NewProxyListFromStrings(
		proxy_list.ProxyListDTO{
			//TODO create local socks5 server for tests
			Addresses: []string{
				"socks5://127.0.0.1:9151", // invalid
				"socks5://127.0.0.1:9150", // valid
			},
			AddressGetter: proxy_list.CreateCircleAddress(0),
		},
	)
	prxyGetter := proxy_list.NewSMTPDialer(proxy_list.NewProxyDialer(proxyList), "")
	emailFrom := ev_email.EmailFromString(DefaultEmail)
	emailTest := ev_email.EmailFromString(emailString)
	mxs, err := net.LookupMX(emailTest.Domain())
	assert.Nil(t, err)

	emptyError := make([]error, 0)

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
		{
			name: emailString,
			fields: fields{
				GetConn:   prxyGetter.Dial, // DirectDial,
				Auth:      nil,
				SendMail:  NewSendMail(),
				FromEmail: emailFrom,
			},
			args: args{
				mxs:   mxs,
				email: emailTest,
			},
			wantErrs: emptyError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := checker{
				dialFunc:  tt.fields.GetConn,
				Auth:      tt.fields.Auth,
				sendMail:  tt.fields.SendMail,
				fromEmail: tt.fields.FromEmail,
			}
			gotErrs := c.Validate(tt.args.mxs, tt.args.email)
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
