package evsmtp

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/proxifier"
	"github.com/stretchr/testify/assert"
	"net"
	"net/smtp"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func mx(domain string, t *testing.T) MXs {
	mxs, err := LookupMX(domain)
	assert.Nil(t, err)

	return mxs
}

func dialFunc(client interface{}, err error) DialFunc {
	return func(addr string) (interface{}, error) {
		return client, err
	}
}

var (
	simpleError  = errors.New("simpleError")
	randomError  = errors.New("randomError")
	mxs          = MXs{&net.MX{}}
	localName    = "localName"
	simpleClient = "simpleClient"
	emailFromStr = "email@from.com"
	emailFrom    = evmail.FromString(emailFromStr)
	emailToStr   = "email@to.com"
	emailTo      = evmail.FromString(emailToStr)
	rAddr        = randomAddress(emailTo)
)

func randomAddress(email evmail.Address) evmail.Address {
	return evmail.FromString("random@" + email.Domain())
}

func mockRandomEmail(t *testing.T, email evmail.Address, err error) RandomEmail {
	return func(domain string) (evmail.Address, error) {
		if domain != email.Domain() {
			t.Errorf("domain of random email is not equal")
		}

		return email, err
	}
}

func Test_checker_Validate(t *testing.T) {
	type fields struct {
		dialFunc    DialFunc
		sendMail    SendMail
		fromEmail   evmail.Address
		localName   string
		randomEmail RandomEmail
	}
	type args struct {
		mx    MXs
		email evmail.Address
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
		{
			name:   "empty mx",
			fields: fields{},
			args:   args{},
			wantErrs: utils.Errs(
				NewError(ConnectionStage, nil),
			),
		},
		{
			name: "cannot connection to mx",
			fields: fields{
				dialFunc: dialFunc(nil, simpleError),
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(
				NewError(ConnectionStage,
					fmt.Errorf(ErrConnection, simpleError),
				),
			),
		},
		{
			name: "Bad hello with localName",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t:    t,
					want: failWant(&sendMailWant{stage: smHello, message: smHello + localName, ret: simpleError}, true),
				},
				localName: localName,
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(
				NewError(HelloStage, simpleError),
			),
		},
		{
			name: "Bad auth",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: failWant(&sendMailWant{
						stage:   smAuth,
						message: smAuth,
						ret:     []interface{}{nil, simpleError},
					}, true),
				},
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(
				NewError(AuthStage, simpleError),
			),
		},
		{
			name: "Bad Mail stage",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: failWant(&sendMailWant{
						stage:   smMail,
						message: smMail + emailFrom.String(),
						ret:     simpleError,
					}, true),
				},
				fromEmail: emailFrom,
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(
				NewError(MailStage, simpleError),
			),
		},
		{
			name: "Problem with generation Random email",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: append(failWant(&sendMailWant{
						stage:   smMail,
						message: smMail + emailFrom.String(),
						ret:     nil,
					}, false),
						sendMailWant{
							stage:   smRcpt,
							message: smRcpt + emailTo.String(),
							ret:     simpleError,
						},
						quitStageWant,
						closeStageWant,
					),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, rAddr, randomError),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RCPTStage, simpleError),
			),
		},
		{
			name: "Problem with RCPT Random email",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: append(failWant(&sendMailWant{
						stage:   smRcpt,
						message: smRcpt + rAddr.String(),
						ret:     simpleError,
					}, false),
						sendMailWant{
							stage:   smRcpt,
							message: smRcpt + emailTo.String(),
							ret:     simpleError,
						},
						quitStageWant,
						closeStageWant,
					),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, rAddr, nil),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RandomRCPTStage, simpleError),
				NewError(RCPTStage, simpleError),
			),
		},
		{
			name: "Quit problem",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: failWant(&sendMailWant{
						stage:   smQuit,
						message: smQuit,
						ret:     simpleError,
					}, true),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, rAddr, nil),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(QuitStage, simpleError),
			),
		},
		{
			name: "Success",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t:    t,
					want: failWant(nil, true),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, rAddr, nil),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: []error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChecker(CheckerDTO{
				DialFunc:    tt.fields.dialFunc,
				SendMail:    tt.fields.sendMail,
				FromEmail:   tt.fields.fromEmail,
				LocalName:   tt.fields.localName,
				RandomEmail: tt.fields.randomEmail,
			})
			if gotErrs := c.Validate(tt.args.mx, tt.args.email); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestChecker_Validate_WithProxy(t *testing.T) {
	return
	// TODO create local socks5 server for tests
	evtests.FunctionalSkip(t)

	type fields struct {
		GetConn   DialFunc
		Auth      smtp.Auth
		SendMail  SendMail
		FromEmail evmail.Address
	}
	type args struct {
		mxs   MXs
		email evmail.Address
	}

	// emailString := "y-numata@senko.ed.jp"
	emailString := "asd@tradepro.net"

	proxyList, _ := proxifier.NewListFromStrings(
		proxifier.ListDTO{
			// TODO create local socks5 server for tests
			Addresses: []string{
				"socks5://127.0.0.1:9151", // invalid
				"socks5://127.0.0.1:9150", // valid
			},
			AddressGetter: proxifier.CreateCircleAddress(0),
		},
	)
	prxyGetter := proxifier.NewSMTPDialer(proxifier.NewProxyDialer(proxyList), "")
	emailFrom := evmail.FromString(DefaultEmail)
	emailTest := evmail.FromString(emailString)
	mxs := mx(emailTest.Domain(), t)

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
