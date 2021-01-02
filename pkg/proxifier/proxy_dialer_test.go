package proxifier

import (
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"net"
	"net/smtp"
	"reflect"
	"testing"
)

const (
	WantGetAddress = "GetAddress"
	WantBan        = "Ban "
	networkTCP     = "tcp"
	networkUDP     = "udp"
)

var (
	tcpConn         = &net.TCPConn{}
	udpConn         = &net.UDPConn{}
	smtpClient      = &smtp.Client{}
	emptySmtpClient = new(smtp.Client)
)

type mockListWant struct {
	value interface{}
	ret   interface{}
}

type mockList struct {
	t    *testing.T
	i    int
	want []mockListWant
}

func (l *mockList) GetAddress() (string, error) {
	want := l.do(WantGetAddress).([]interface{})
	return want[0].(string), evtests.ToError(want[1])
}

func (l *mockList) Ban(s string) bool {
	return l.do(WantBan + s).(bool)
}

func (l *mockList) do(value interface{}) interface{} {
	if l.i >= len(l.want) {
		l.t.Fatalf("Invalid command %q", value)
	}

	if value != l.want[l.i].value {
		l.t.Fatalf("Invalid command, got %q, want %q", value, l.want[l.i].value)
	}
	l.i++

	return l.want[l.i-1].ret
}

type dialFuncWant struct {
	network string
	addr    string
	conn    net.Conn
	err     error
}

type mockDialer struct {
	t    *testing.T
	i    int
	want []dialFuncWant
}

func (d *mockDialer) Dial(network, addr string) (c net.Conn, err error) {
	if d.i >= len(d.want) {
		d.t.Fatalf("Invalid arguments %v, %v", network, addr)
	}

	if network != d.want[d.i].network {
		d.t.Fatalf("Invalid network %q, want %q", network, d.want[d.i].network)
	}
	if addr != d.want[d.i].addr {
		d.t.Fatalf("Invalid addr %q, want %q", addr, d.want[d.i].addr)
	}

	d.i++

	return d.want[d.i-1].conn, d.want[d.i-1].err
}

type dialProxyFuncWant struct {
	dialFuncWant
	proxyURI string
}

func mockProxyDialFunc(t *testing.T, want []dialProxyFuncWant) ProxyDialerFunc {
	i := 0
	return func(proxyURI string) func(string, string) (net.Conn, error) {
		if i >= len(want) {
			t.Fatalf("Invalid command %q", proxyURI)
		}

		if proxyURI != want[i].proxyURI {
			t.Fatalf("Invalid command, got %q, want %q", proxyURI, want[i].addr)
		}

		return func(network string, addr string) (net.Conn, error) {
			if i >= len(want) {
				t.Fatalf("Invalid command %q", proxyURI)
			}

			if network != want[i].network {
				t.Fatalf("Invalid command, got %q, want %q", network, want[i].network)
			}
			if addr != want[i].addr {
				t.Fatalf("Invalid command, got %q, want %q", addr, want[i].addr)
			}

			i++
			return want[i-1].conn, want[i-1].err
		}
	}
}

func Test_dialer_Dial(t *testing.T) {
	type fields struct {
		list       List
		dialerFunc ProxyDialerFunc
	}
	type args struct {
		network string
		addr    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantC   net.Conn
		wantErr bool
	}{
		{
			name: "successful the first time",
			fields: fields{
				list: &mockList{
					t: t,
					want: []mockListWant{
						{
							value: WantGetAddress,
							ret:   []interface{}{addressFirstWithPort, nil},
						},
					},
				},
				dialerFunc: mockProxyDialFunc(t, []dialProxyFuncWant{
					{
						proxyURI: addressFirstWithPort,
						dialFuncWant: dialFuncWant{
							network: networkTCP,
							addr:    addressSecond,
							conn:    tcpConn,
							err:     nil,
						},
					},
				}),
			},
			args: args{
				network: networkTCP,
				addr:    addressSecond,
			},
			wantC:   tcpConn,
			wantErr: false,
		},
		{
			name: "successful the second time",
			fields: fields{
				list: &mockList{
					t: t,
					want: []mockListWant{
						{
							value: WantGetAddress,
							ret:   []interface{}{addressFirstWithPort, nil},
						},
						{
							value: WantBan + addressFirstWithPort,
							ret:   true,
						},
						{
							value: WantGetAddress,
							ret:   []interface{}{addressSecond, nil},
						},
					},
				},
				dialerFunc: mockProxyDialFunc(t, []dialProxyFuncWant{
					{
						proxyURI: addressFirstWithPort,
						dialFuncWant: dialFuncWant{
							network: networkUDP,
							addr:    addressThird,
							conn:    nil,
							err:     simpleError,
						},
					},
					{
						proxyURI: addressSecond,
						dialFuncWant: dialFuncWant{
							network: networkUDP,
							addr:    addressThird,
							conn:    tcpConn,
							err:     nil,
						},
					},
				}),
			},
			args: args{
				network: networkUDP,
				addr:    addressThird,
			},
			wantC:   tcpConn,
			wantErr: false,
		},
		{
			name: "not enough proxies",
			fields: fields{
				list: &mockList{
					t: t,
					want: []mockListWant{
						{
							value: WantGetAddress,
							ret:   []interface{}{EmptyAddress, simpleError},
						},
					},
				},
				dialerFunc: nil,
			},
			args: args{
				network: networkTCP,
				addr:    addressFirstWithPort,
			},
			wantC:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dialer{
				list:       tt.fields.list,
				dialerFunc: tt.fields.dialerFunc,
			}
			gotC, err := d.Dial(tt.args.network, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("Dial() gotC = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestNewProxyDialer(t *testing.T) {
	type args struct {
		list       List
		dialerFunc ProxyDialerFunc
	}
	tests := []struct {
		name string
		args args
		want proxy.Dialer
	}{
		{
			name: "empty dialerFunc",
			args: args{
				list:       &mockList{},
				dialerFunc: nil,
			},
			want: &dialer{
				list:       &mockList{},
				dialerFunc: socks.Dial,
			},
		},
		{
			name: "filled",
			args: args{
				list:       &mockList{},
				dialerFunc: mockProxyDialFunc(t, nil),
			},
			want: &dialer{
				list:       &mockList{},
				dialerFunc: mockProxyDialFunc(t, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProxyDialer(tt.args.list, tt.args.dialerFunc); !reflect.DeepEqual(got.(*dialer).list, tt.want.(*dialer).list) &&
				fmt.Sprintf("%v", got.(*dialer).dialerFunc) != fmt.Sprintf("%v", tt.want.(*dialer).dialerFunc) {
				t.Errorf("NewProxyDialer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_smtpDialer_Dial(t *testing.T) {
	type fields struct {
		dialer        proxy.Dialer
		network       string
		smtpNewClient func(conn net.Conn, host string) (*smtp.Client, error)
	}
	type args struct {
		addr string
	}

	defaultMockDialer := &mockDialer{
		t: t,
		want: []dialFuncWant{
			{
				network: TCPConnection,
				addr:    addressFirstWithPort,
				conn:    tcpConn,
				err:     nil,
			},
		},
	}

	errorMockDialer := &mockDialer{
		t: t,
		want: []dialFuncWant{
			{
				network: UDPConnection,
				addr:    addressFirstWithPort,
				conn:    tcpConn,
				err:     simpleError,
			},
		},
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		want       interface{}
		wantErr    bool
		wantDialer *smtpDialer
	}{
		{
			name: "default " + addressFirstWithPort,
			fields: fields{
				dialer:  defaultMockDialer,
				network: "",
				smtpNewClient: func(conn net.Conn, host string) (*smtp.Client, error) {
					require.Equal(t, tcpConn, conn)
					require.Equal(t, addressFirst, host)

					return smtpClient, nil
				},
			},
			args: args{
				addr: addressFirstWithPort,
			},
			want:    smtpClient,
			wantErr: false,
			wantDialer: &smtpDialer{
				dialer:  defaultMockDialer,
				network: TCPConnection,
			},
		},
		{
			name: "problem proxy " + addressFirstWithPort,
			fields: fields{
				dialer:  errorMockDialer,
				network: UDPConnection,
			},
			args: args{
				addr: addressFirstWithPort,
			},
			want:    nil,
			wantErr: true,
			wantDialer: &smtpDialer{
				dialer:  errorMockDialer,
				network: UDPConnection,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSMTPDialer(tt.fields.dialer, tt.fields.network)
			if !reflect.DeepEqual(p, tt.wantDialer) {
				t.Errorf("smtpDialer() p = %v, wantDialer %v", p, tt.wantDialer)
				return
			}

			smtpNewClient = tt.fields.smtpNewClient

			got, err := p.Dial(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dial() got = %v, want %v", got, tt.want)
			}
		})
	}
	smtpNewClient = smtp.NewClient
}
