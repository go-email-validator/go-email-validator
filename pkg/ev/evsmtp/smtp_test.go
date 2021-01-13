package evsmtp

import (
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/eko/gocache/marshaler"
	"github.com/eko/gocache/store"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtp_client"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/proxifier"
	mock_evcache "github.com/go-email-validator/go-email-validator/test/mock/ev/evcache"
	mock_evmail "github.com/go-email-validator/go-email-validator/test/mock/ev/evmail"
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

func dialFunc(client smtp_client.SMTPClient, err error) DialFunc {
	return func(addr string) (smtp_client.SMTPClient, error) {
		return client, err
	}
}

var (
	simpleError    = errors.New("simpleError")
	randomError    = errors.New("randomError")
	mxs            = MXs{&net.MX{Host: "127.0.0.1"}}
	localName      = "localName"
	emptyLocalName = ""
	simpleClient   = &smtp.Client{}
	emailFromStr   = "email@from.com"
	emailFrom      = evmail.FromString(emailFromStr)
	emailToStr     = "email@to.com"
	emailTo        = evmail.FromString(emailToStr)
	randomAddress  = getRandomAddress(emailTo)
	validEmail     = mock_evmail.GetValidTestEmail()
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

func getLocalIP() string {
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
				ErrConnection,
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
					simpleError,
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
							stage:   smRCPTs,
							message: smRCPTs + emailTo.String(),
							ret:     simpleError,
						},
						quitStageWant,
						closeStageWant,
					),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, randomAddress, randomError),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RandomRCPTStage, randomError),
				NewError(RCPTsStage, simpleError),
			),
		},
		{
			name: "Problem with RCPTs Random email",
			fields: fields{
				dialFunc: dialFunc(simpleClient, nil),
				sendMail: &mockSendMail{
					t: t,
					want: append(failWant(&sendMailWant{
						stage:   smRCPTs,
						message: smRCPTs + randomAddress.String(),
						ret:     simpleError,
					}, false),
						sendMailWant{
							stage:   smRCPTs,
							message: smRCPTs + emailTo.String(),
							ret:     simpleError,
						},
						quitStageWant,
						closeStageWant,
					),
				},
				fromEmail:   emailFrom,
				randomEmail: mockRandomEmail(t, randomAddress, nil),
			},
			args: args{
				mx:    mxs,
				email: emailTo,
			},
			wantErrs: utils.Errs(
				NewError(RandomRCPTStage, simpleError),
				NewError(RCPTsStage, simpleError),
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
				randomEmail: mockRandomEmail(t, randomAddress, nil),
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
				randomEmail: mockRandomEmail(t, randomAddress, nil),
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
			gotErrs := c.Validate(tt.args.mx, tt.args.email)
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
		"EHLO localhost",
		"HELO localhost",
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

	localIp := getLocalIP()

	invalidProxies := []string{
		"socks5://0.0.0.0:0", //invalid
	}

	type fields struct {
		GetConn     DialFunc
		Auth        smtp.Auth
		SendMail    SendMail
		FromEmail   evmail.Address
		Localhost   string
		RandomEmail RandomEmail
		Port        int
		Server      []string
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
				GetConn:     Dial,
				Auth:        nil,
				SendMail:    NewSendMail(nil),
				FromEmail:   emailFrom,
				Localhost:   "localhost",
				RandomEmail: mockRandomEmail(t, getRandomAddress(emailTest), nil),
				Server:      successServer,
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
		{
			name: "with proxy success after ban",
			fields: fields{
				GetConn:     getSMTPProxy(nil, append(invalidProxies, proxyList...)...).Dial,
				Auth:        nil,
				SendMail:    NewSendMail(nil),
				FromEmail:   emailFrom,
				Localhost:   "localhost",
				RandomEmail: mockRandomEmail(t, getRandomAddress(emailTest), nil),
				Server:      successServer,
			},
			args: args{
				mxs: MXs{&net.MX{
					Host: localIp,
				}},
				email: emailTest,
			},
			wantErrs: []error{NewError(RandomRCPTStage, &textproto.Error{
				Code: 550,
				Msg:  "address does not exist",
			})},
			wantSMTP: successWantSMTP,
		},
		{
			name: "with invalid proxy",
			fields: fields{
				GetConn:     getSMTPProxy(nil, invalidProxies...).Dial,
				Auth:        nil,
				SendMail:    NewSendMail(nil),
				FromEmail:   emailFrom,
				Localhost:   "localhost",
				RandomEmail: mockRandomEmail(t, getRandomAddress(emailTest), nil),
			},
			args: args{
				mxs:   mxs,
				email: emailTest,
			},
			wantErrs: []error{NewError(ConnectionStage, proxifier.ErrEmptyPool)},
			wantSMTP: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, done := mock_evsmtp.Server(t, tt.fields.Server, time.Second)

			if tt.fields.Port == 0 {
				u, _ := url.Parse("http://" + addr)
				tt.fields.Port, _ = strconv.Atoi(u.Port())
			}

			c := checker{
				dialFunc:    tt.fields.GetConn,
				Auth:        tt.fields.Auth,
				sendMail:    tt.fields.SendMail,
				fromEmail:   tt.fields.FromEmail,
				localName:   tt.fields.Localhost,
				randomEmail: tt.fields.RandomEmail,
				port:        tt.fields.Port,
			}
			c.RandomRCPT = &ARandomRCPT{fn: c.randomRCPT}

			gotErrs := c.Validate(tt.args.mxs, tt.args.email)
			actualClient := <-done

			wantSMTP := strings.Join(tt.wantSMTP, mock_evsmtp.Separator)
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

	errs := []error{simpleError}
	errsAlias := []AliasError{simpleError}
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
					mock := mock_evcache.NewMockInterface(ctrl)
					mock.EXPECT().Get(validEmail.Domain()).Return(errsAlias, nil).Times(1)

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
					mock.EXPECT().Call(validEmail).Return(errs).Times(1)

					return mock
				},
				cache: func() evcache.Interface {
					mock := mock_evcache.NewMockInterface(ctrl)
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
			if gotErrs := c.RandomRCPT(tt.args.email); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
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
				email: mock_evmail.GetValidTestEmail(),
			},
			want: mock_evmail.GetValidTestEmail().Domain(),
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

	err = cache.Set(key, ConvertErrorsToEVSMTPErrors(cacheErrs))
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
					mock.EXPECT().Call(validEmail).Return(cacheErrs).Times(1)

					return mock
				},
				cache: func() evcache.Interface {
					bigCacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
					require.Nil(t, err)
					bigCacheStore := store.NewBigcache(bigCacheClient, nil)

					marshal := marshaler.New(bigCacheStore)

					// Add value to cache
					err = marshal.Set(DefaultRandomCacheKeyGetter(validEmail), cacheErrs, nil)
					require.Nil(t, err)

					return evcache.NewCacheMarshaller(marshal, func() interface{} {
						return []error{}
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
			if gotErrs := c.RandomRCPT(tt.args.email); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("RandomRCPT() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}
