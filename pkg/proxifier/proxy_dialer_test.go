package proxifier

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"net"
	"reflect"
	"testing"
)

const (
	WantGetAddress = "GetAddress"
	WantBan        = "Ban "
	networkTCP     = "tcp"
	networkUDP     = "udp"
)

var tcpConn = &net.TCPConn{}

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
	proxyURI string
	network  string
	addr     string
	conn     net.Conn
	err      error
}

func mockDialFunc(t *testing.T, want []dialFuncWant) ProxyDialerFunc {
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
							ret:   []interface{}{addressFirst, nil},
						},
					},
				},
				dialerFunc: mockDialFunc(t, []dialFuncWant{
					{
						proxyURI: addressFirst,
						network:  networkTCP,
						addr:     addressSecond,
						conn:     tcpConn,
						err:      nil,
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
							ret:   []interface{}{addressFirst, nil},
						},
						{
							value: WantBan + addressFirst,
							ret:   true,
						},
						{
							value: WantGetAddress,
							ret:   []interface{}{addressSecond, nil},
						},
					},
				},
				dialerFunc: mockDialFunc(t, []dialFuncWant{
					{
						proxyURI: addressFirst,
						network:  networkUDP,
						addr:     addressThird,
						conn:     nil,
						err:      simpleError,
					},
					{
						proxyURI: addressSecond,
						network:  networkUDP,
						addr:     addressThird,
						conn:     tcpConn,
						err:      nil,
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
				addr:    addressFirst,
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
