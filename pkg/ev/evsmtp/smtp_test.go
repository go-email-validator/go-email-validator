package evsmtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/eko/gocache/marshaler"
	"github.com/eko/gocache/store"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/proxifier"
	mockevcache "github.com/go-email-validator/go-email-validator/test/mock/ev/evcache"
	mockevmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
	"github.com/go-email-validator/go-email-validator/test/mock/ev/evsmtp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net"
	"net/smtp"
	"net/textproto"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

const EnvPath = "../../../.env"

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func dialFunc(t *testing.T, client smtpclient.SMTPClient, err error, wantCtx context.Context, wantAddr, wantProxy string, sleep time.Duration) DialFunc {
	return func(ctx context.Context, addr, proxy string) (smtpclient.SMTPClient, error) {
		require.Equal(t, utils.StructName(wantCtx), utils.StructName(ctx))
		require.Equal(t, addr, wantAddr)
		require.Equal(t, wantProxy, proxy)

		time.Sleep(sleep)

		return client, err
	}
}

var (
	localhost      = "127.0.0.1"
	smtpLocalhost  = localhost + ":25"
	errorSimple    = errors.New("errorSimple")
	errorRandom    = errors.New("errorRandom")
	mxs            = MXs{&net.MX{Host: localhost}}
	emptyLocalName = ""
	simpleClient   = &smtp.Client{}
	emailFromStr   = "email@from.com"
	emailFrom      = evmail.FromString(emailFromStr)
	emailToStr     = "email@to.com"
	emailTo        = evmail.FromString(emailToStr)
	randomAddress  = getRandomAddress(emailTo)
	validEmail     = mockevmail.GetValidTestEmail()
	getMockKey     = func(t *testing.T, wantEmail evmail.Address, ret interface{}) func(email evmail.Address) interface{} {
		return func(email evmail.Address) interface{} {
			require.Equal(t, wantEmail, email)
			return ret
		}
	}
)

func getRandomAddress(email evmail.Address) evmail.Address {
	return evmail.FromString("random.which.did.not.exist@" + email.Domain())
}

func mockRandomEmail(t *testing.T, email evmail.Address, err error) RandomEmail {
	return func(domain string) (evmail.Address, error) {
		if domain != email.Domain() {
			t.Errorf("domain of random email is not equal")
		}

		return email, err
	}
}

func getSMTPProxy(dialerFunc proxifier.ProxyDialerFunc, proxies ...string) proxifier.SMTPDialler {
	proxyList, _ := proxifier.NewListFromStrings(
		proxifier.ListDTO{
			Addresses:     proxies,
			AddressGetter: proxifier.CreateCircleAddress(0),
		},
	)
	return proxifier.NewSMTPDialer(proxifier.NewProxyDialer(proxyList, dialerFunc), "")
}

// TODO delete after remove proxifier
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func Test_checker_Validate(t *testing.T) {
	type fields struct {
		sendMailFactory SendMailDialerFactory
		randomEmail     RandomEmail
		options         Options
	}
	type args struct {
		mx    MXs
		email evmail.Address
	}

	successDialFunc := dialFunc(t, simpleClient, nil, context.Background(), smtpLocalhost, "", 0)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	errConnection := NewError(ConnectionStage, errors.New(ErrConnectionMsg))

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
		{
			name:     "empty mx",
			args:     args{},
			wantErrs: utils.Errs(errConnection),
		},
		{
			name: "cannot connection to mx",
			fields: fields{
				sendMailFactory: NewSendMailFactory(dialFunc(t, nil, errorSimple, context.Background(), smtpLocalhost, "", 0), nil),
				options:         &options{},
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(errConnection),
		},
		{
			name: "Bad hello with helloName",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t:    t,
							want: failWant(&sendMailWant{stage: smHello, message: smHello + helloName, ret: errorSimple}, true),
						}
					}),
				options: &options{
					helloName: helloName,
				},
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(NewError(HelloStage, errorSimple)),
		},
		{
			name: "Bad auth",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: failWant(&sendMailWant{
								stage:   smAuth,
								message: smAuth,
								ret:     []interface{}{nil, errorSimple},
							}, true),
						}
					}),
				options: &options{},
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(NewError(AuthStage, errorSimple)),
		},
		{
			name: "Bad Mail stage",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: failWant(&sendMailWant{
								stage:   smMail,
								message: smMail + emailFrom.String(),
								ret:     errorSimple,
							}, true),
						}
					}),
				options: &options{
					emailFrom: emailFrom,
				},
			},
			args: args{
				mx: mxs,
			},
			wantErrs: utils.Errs(NewError(MailStage, errorSimple)),
		},
		{
			name: "Problem with generation Random email",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: append(failWant(&sendMailWant{
								stage:   smMail,
								message: smMail + emailFrom.String(),
								ret:     nil,
							}, false),
								sendMailWant{
									stage:   smRCPTs,
									message: smRCPTs + emailTo.String(),
									ret:     errorSimple,
								},
								quitStageWant,
								closeStageWant,
							),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, errorRandom),
				options: &options{
					emailFrom: emailFrom,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RandomRCPTStage, errorRandom),
				NewError(RCPTsStage, errorSimple),
			),
		},
		{
			name: "Problem with RCPTs Random email",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: append(failWant(&sendMailWant{
								stage:   smRCPTs,
								message: smRCPTs + randomAddress.String(),
								ret:     errorSimple,
							}, false),
								sendMailWant{
									stage:   smRCPTs,
									message: smRCPTs + emailTo.String(),
									ret:     errorSimple,
								},
								quitStageWant,
								closeStageWant,
							),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom: emailFrom,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RandomRCPTStage, errorSimple),
				NewError(RCPTsStage, errorSimple),
			),
		},
		{
			name: "Quit problem",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: failWant(&sendMailWant{
								stage:   smQuit,
								message: smQuit,
								ret:     errorSimple,
							}, true),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom: emailFrom,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(NewError(QuitStage, errorSimple)),
		},
		{
			name: "Success",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t:    t,
							want: failWant(nil, true),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom: emailFrom,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: []error{},
		},
		{
			name: "with timeout success",
			fields: fields{
				sendMailFactory: NewSendMailCustom(
					dialFunc(t, simpleClient, nil, ctxTimeout, smtpLocalhost, "", 0),
					nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t:    t,
							want: failWant(nil, true),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom:   emailFrom,
					timeoutCon:  5 * time.Second,
					timeoutResp: 5 * time.Second,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: []error{},
		},
		{
			name: "with expired connection timeout",
			fields: fields{
				sendMailFactory: NewSendMailCustom(
					dialFunc(t, simpleClient, nil, ctxTimeout, smtpLocalhost, "", 2*time.Millisecond),
					nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t:    t,
							want: failWant(nil, true),
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom:  emailFrom,
					timeoutCon: 1,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(errConnection),
		},
		{
			name: "with expired response timeout",
			fields: fields{
				sendMailFactory: NewSendMailCustom(successDialFunc, nil,
					func(client smtpclient.SMTPClient, tlsConfig *tls.Config) SendMail {
						return &mockSendMail{
							t: t,
							want: []sendMailWant{
								{
									sleep:   2 * time.Millisecond,
									stage:   smHello,
									message: smHelloLocalhost,
									ret:     context.DeadlineExceeded,
								},
								closeStageWant,
							},
						}
					}),
				randomEmail: mockRandomEmail(t, randomAddress, nil),
				options: &options{
					emailFrom:   emailFrom,
					timeoutResp: 1 * time.Millisecond,
				},
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(NewError(HelloStage, context.DeadlineExceeded)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChecker(CheckerDTO{
				SendMailFactory: tt.fields.sendMailFactory,
				RandomEmail:     tt.fields.randomEmail,
				Options:         tt.fields.options,
			})
			gotErrs := c.Validate(tt.args.mx, NewInput(tt.args.email, nil))
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestChecker_Validate_WithProxy_Local(t *testing.T) {
	evtests.FunctionalSkip(t)

	successServer := []string{
		"220 hello world",
		"502 EH?",
		"250 mx.google.com at your service",
		"250 Sender ok",
		"550 address does not exist",
		"250 Receiver ok",
		"221 Goodbye",
	}
	successWantSMTP := []string{
		"EHLO helloName",
		"HELO helloName",
		"MAIL FROM:<user@example.org>",
		"RCPT TO:<random.which.did.not.exist@tradepro.net>",
		"RCPT TO:<asd@tradepro.net>",
		"QUIT",
		"",
	}

	utils.LoadEnv(EnvPath)
	proxyList := proxifier.EnvProxies()
	if len(proxyList) == 0 {
		t.Error("PROXIES env should be set")
		return
	}

	// TODO delete after deleting of proxifier
	//lIP := localIP()
	//
	//invalidProxies := []string{
	//	"socks5://0.0.0.0:0", //invalid
	//}

	type fields struct {
		SendMailFactory SendMailDialerFactory
		Auth            smtp.Auth
		RandomEmail     RandomEmail
		Server          []string
		OptionsDTO      OptionsDTO
	}
	type args struct {
		mxs   MXs
		email evmail.Address
	}

	emailString := "asd@tradepro.net"

	emailFrom := evmail.FromString(DefaultEmail)
	emailTest := evmail.FromString(emailString)

	emptyError := make([]error, 0)
	_ = emptyError

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
		wantSMTP []string
	}{
		{
			name: "without proxy",
			fields: fields{
				SendMailFactory: NewSendMailFactory(DirectDial, nil),
				Auth:            nil,
				RandomEmail:     mockRandomEmail(t, getRandomAddress(emailTest), nil),
				Server:          successServer,
				OptionsDTO: OptionsDTO{
					EmailFrom: emailFrom,
					HelloName: helloName,
				},
			},
			args: args{
				mxs:   mxs,
				email: emailTest,
			},
			wantErrs: []error{NewError(RandomRCPTStage, &textproto.Error{
				Code: 550,
				Msg:  "address does not exist",
			})},
			wantSMTP: successWantSMTP,
		},
		// TODO delete after deleting of proxifier
		//{
		//	name: "with proxy success after ban",
		//	fields: fields{
		//		GetConn:     getSMTPProxy(nil, append(invalidProxies, proxyList...)...).DialContext,
		//		Auth:        nil,
		//		SendMail:    NewSendMail(nil),
		//		RandomEmail: mockRandomEmail(t, getRandomAddress(emailTest), nil),
		//		Server:      successServer,
		//		OptionsDTO: OptionsDTO{
		//			EmailFrom: emailFrom,
		//			HelloName: helloName,
		//		},
		//	},
		//	args: args{
		//		mxs: MXs{&net.MX{
		//			Host: lIP,
		//		}},
		//		email: emailTest,
		//	},
		//	wantErrs: []error{NewError(RandomRCPTStage, &textproto.Error{
		//		Code: 550,
		//		Msg:  "address does not exist",
		//	})},
		//	wantSMTP: successWantSMTP,
		//},
		//{
		//	name: "with invalid proxy",
		//	fields: fields{
		//		GetConn:     getSMTPProxy(nil, invalidProxies...).DialContext,
		//		Auth:        nil,
		//		SendMail:    NewSendMail(nil),
		//		RandomEmail: mockRandomEmail(t, getRandomAddress(emailTest), nil),
		//		OptionsDTO: OptionsDTO{
		//			EmailFrom: emailFrom,
		//			HelloName: helloName,
		//		},
		//	},
		//	args: args{
		//		mxs:   mxs,
		//		email: emailTest,
		//	},
		//	wantErrs: []error{NewError(ConnectionStage, proxifier.ErrEmptyPool)},
		//	wantSMTP: []string{},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, done := mockevsmtp.Server(t, tt.fields.Server, time.Second)

			if tt.fields.OptionsDTO.Port == 0 {
				u, _ := url.Parse("http://" + addr)
				tt.fields.OptionsDTO.Port, _ = strconv.Atoi(u.Port())
			}

			c := checker{
				sendMailFactory: tt.fields.SendMailFactory,
				Auth:            tt.fields.Auth,
				randomEmail:     tt.fields.RandomEmail,
				options:         NewOptions(tt.fields.OptionsDTO),
			}
			c.RandomRCPT = &ARandomRCPT{fn: c.randomRCPT}

			gotErrs := c.Validate(tt.args.mxs, NewInput(tt.args.email, nil))
			actualClient := <-done

			wantSMTP := strings.Join(tt.wantSMTP, mockevsmtp.Separator)
			if wantSMTP != actualClient {
				t.Errorf("Got:\n%s\nExpected:\n%s", actualClient, wantSMTP)
			}

			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("Validate() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_checkerCacheRandomRCPT_RandomRCPT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		checkerWithRandomRPCT func() CheckerWithRandomRCPT
		cache                 func() evcache.Interface
		getKey                RandomCacheKeyGetter
	}
	type args struct {
		email evmail.Address
	}

	errs := []error{errorSimple}
	errsAlias := []AliasError{errorSimple}
	emptyChecker := func() CheckerWithRandomRCPT {
		mock := NewMockCheckerWithRandomRCPT(ctrl)
		mock.EXPECT().get().Return(nil).Times(1)
		mock.EXPECT().set(gomock.Any()).Times(1)

		return mock
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
		{
			name: "with cache",
			fields: fields{
				checkerWithRandomRPCT: emptyChecker,
				cache: func() evcache.Interface {
					mock := mockevcache.NewMockInterface(ctrl)
					mock.EXPECT().Get(validEmail.Domain()).Return(&errs, nil).Times(1)

					return mock
				},
				getKey: getMockKey(t, validEmail, validEmail.Domain()),
			},
			args: args{
				email: validEmail,
			},
			wantErrs: errs,
		},
		{
			name: "without cache",
			fields: fields{
				checkerWithRandomRPCT: func() CheckerWithRandomRCPT {
					mock := NewMockCheckerWithRandomRCPT(ctrl)
					mock.EXPECT().get().Return(mock.Call).Times(1)
					mock.EXPECT().set(gomock.Any()).Times(1)
					mock.EXPECT().Call(gomock.Any(), validEmail).Return(errs).Times(1)

					return mock
				},
				cache: func() evcache.Interface {
					mock := mockevcache.NewMockInterface(ctrl)
					mock.EXPECT().Get(validEmail.Domain()).Return(nil, nil).Times(1)
					mock.EXPECT().Set(validEmail.Domain(), errsAlias).Return(nil).Times(1)

					return mock
				},
				getKey: getMockKey(t, validEmail, validEmail.Domain()),
			},
			args: args{
				email: validEmail,
			},
			wantErrs: errs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCheckerCacheRandomRCPT(tt.fields.checkerWithRandomRPCT(), tt.fields.cache(), tt.fields.getKey).(*checkerCacheRandomRCPT)
			if gotErrs := c.RandomRCPT(nil, tt.args.email); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("RandomRCPT() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func TestDefaultRandomCacheKeyGetter(t *testing.T) {
	type args struct {
		email evmail.Address
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "success",
			args: args{
				email: mockevmail.GetValidTestEmail(),
			},
			want: mockevmail.GetValidTestEmail().Domain(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultRandomCacheKeyGetter(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultRandomCacheKeyGetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCheckerCacheRandomRCPT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		checker func() CheckerWithRandomRCPT
		cache   evcache.Interface
		getKey  RandomCacheKeyGetter
	}
	tests := []struct {
		name string
		args args
		want Checker
	}{
		{
			name: "fill empty",
			args: args{
				checker: func() CheckerWithRandomRCPT {
					mock := NewMockCheckerWithRandomRCPT(ctrl)
					mock.EXPECT().get().Return(nil).Times(1)
					mock.EXPECT().set(gomock.Any()).Times(1)

					return mock
				},
				cache:  nil,
				getKey: nil,
			},
			want: &checkerCacheRandomRCPT{
				getKey:     DefaultRandomCacheKeyGetter,
				randomRCPT: &ARandomRCPT{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCheckerCacheRandomRCPT(tt.args.checker(), tt.args.cache, tt.args.getKey)

			gotChecker := got.(*checkerCacheRandomRCPT)
			gotGetKey := gotChecker.getKey
			gotChecker.getKey = nil
			gotChecker.CheckerWithRandomRCPT = nil
			want := tt.want.(*checkerCacheRandomRCPT)
			wantGetKey := want.getKey
			want.getKey = nil

			if !reflect.DeepEqual(got, tt.want) || fmt.Sprint(gotGetKey) != fmt.Sprint(wantGetKey) {
				t.Errorf(
					"NewCheckerCacheRandomRCPT() = %v, want %v\n gotGetKey = %v, wantGetKey %v",
					got, tt.want, gotGetKey, wantGetKey)
			}
		})
	}
}

var cacheErrs = []error{
	NewError(1, &textproto.Error{Code: 505, Msg: "msg1"}),
	NewError(1, errors.New("msg2")),
}

func Test_Cache(t *testing.T) {
	bigCacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
	require.Nil(t, err)
	bigCacheStore := store.NewBigcache(bigCacheClient, nil)

	marshal := marshaler.New(bigCacheStore)

	cache := evcache.NewCacheMarshaller(marshal, func() interface{} {
		return new([]error)
	}, nil)

	key := "key"

	err = cache.Set(key, ErrorsToEVSMTPErrors(cacheErrs))
	require.Nil(t, err)

	got, err := cache.Get(key)
	require.Nil(t, err)
	require.Equal(t, cacheErrs, *got.(*[]error))
}

func Test_checkerCacheRandomRCPT_RandomRCPT_RealCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		CheckerWithRandomRCPT func() CheckerWithRandomRCPT
		randomRCPT            RandomRCPT
		cache                 func() evcache.Interface
	}
	type args struct {
		email evmail.Address
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
		{
			name: "with cache",
			fields: fields{
				CheckerWithRandomRCPT: func() CheckerWithRandomRCPT {
					mock := NewMockCheckerWithRandomRCPT(ctrl)
					mock.EXPECT().get().Return(mock.Call).Times(1)
					mock.EXPECT().set(gomock.Any()).Times(1)

					return mock
				},
				cache: func() evcache.Interface {
					bigCacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
					require.Nil(t, err)
					bigCacheStore := store.NewBigcache(bigCacheClient, nil)

					marshal := marshaler.New(bigCacheStore)

					// Add value to cache
					key := DefaultRandomCacheKeyGetter(validEmail)
					err = marshal.Set(key, ErrorsToEVSMTPErrors(cacheErrs), nil)
					require.Nil(t, err)

					return evcache.NewCacheMarshaller(marshal, func() interface{} {
						return new([]error)
					}, nil)
				},
			},
			args: args{
				email: validEmail,
			},
			wantErrs: cacheErrs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCheckerCacheRandomRCPT(tt.fields.CheckerWithRandomRCPT(), tt.fields.cache(), DefaultRandomCacheKeyGetter).(*checkerCacheRandomRCPT)
			if gotErrs := c.RandomRCPT(nil, tt.args.email); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("RandomRCPT() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
