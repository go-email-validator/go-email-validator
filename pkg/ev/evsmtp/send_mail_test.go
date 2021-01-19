package evsmtp

import (
	"crypto/tls"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp/smtpclient"
	mocksmtpclient "github.com/go-email-validator/go-email-validator/test/mock/ev/evsmtp/smtpclient"
	"github.com/golang/mock/gomock"
	"io"
	"reflect"
	"testing"
)

func TestNewSendMail(t *testing.T) {
	tlsConfig := &tls.Config{}

	type args struct {
		tlsConfig *tls.Config
	}
	tt := struct {
		name string
		args args
		want SendMail
	}{
		args: args{
			tlsConfig: tlsConfig,
		},
		want: &sendMail{
			tlsConfig: tlsConfig,
		},
	}

	if got := NewSendMail(tt.args.tlsConfig); !reflect.DeepEqual(got, tt.want) {
		t.Errorf("NewSendMail() = %v, want %v", got, tt.want)
	}
}

func Test_sendMail_Client(t *testing.T) {
	type fields struct {
		client    smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "filled",
			fields: fields{
				client: simpleClient,
			},
			want: simpleClient,
		},
		{
			name: "nil",
			fields: fields{
				client: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client,
				tlsConfig: tt.fields.TLSConfig,
			}
			if got := s.Client(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendMail_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Close().Return(nil).Times(1)

					return smtpMock
				},
			},
			want: nil,
		},
		{
			name: "with error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Close().Return(errorSimple).Times(1)

					return smtpMock
				},
			},
			want: errorSimple,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			if err := s.Close(); !reflect.DeepEqual(tt.want, err) {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func Test_sendMail_Data(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	tests := []struct {
		name    string
		fields  fields
		want    io.WriteCloser
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Data().Return(mockWriterInstance, nil).Times(1)

					return smtpMock
				},
			},
			want:    mockWriterInstance,
			wantErr: nil,
		},
		{
			name: "error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Data().Return(nil, errorSimple).Times(1)

					return smtpMock
				},
			},
			want:    nil,
			wantErr: errorSimple,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			got, err := s.Data()
			if !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Data() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendMail_Hello(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	type args struct {
		localName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Hello(helloName).Return(nil).Times(1)

					return smtpMock
				},
			},
			args: args{
				localName: helloName,
			},
			wantErr: nil,
		},
		{
			name: "error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Hello(emptyLocalName).Return(errorSimple).Times(1)

					return smtpMock
				},
			},
			args: args{
				localName: emptyLocalName,
			},
			wantErr: errorSimple,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			if err := s.Hello(tt.args.localName); !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Hello() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sendMail_Mail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	type args struct {
		from string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Mail(emailFromStr).Return(nil).Times(1)

					return smtpMock
				},
			},
			args: args{
				from: emailFromStr,
			},
			wantErr: nil,
		},
		{
			name: "error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Mail(emailFromStr).Return(errorSimple).Times(1)

					return smtpMock
				},
			},
			args: args{
				from: emailFromStr,
			},
			wantErr: errorSimple,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			if err := s.Mail(tt.args.from); !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Mail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sendMail_Quit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Quit().Return(nil).Times(1)

					return smtpMock
				},
			},
			wantErr: nil,
		},
		{
			name: "error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					smtpMock.EXPECT().Quit().Return(errorSimple).Times(1)

					return smtpMock
				},
			},
			wantErr: errorSimple,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			if err := s.Quit(); !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Quit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sendMail_RCPTs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addrs := []string{emailToStr, emailFromStr}
	type fields struct {
		client    func() smtpclient.SMTPClient
		TLSConfig *tls.Config
	}
	type args struct {
		addrs []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]error
	}{
		{
			name: "success",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					firstCall := smtpMock.EXPECT().Rcpt(emailToStr).Return(nil).Times(1)
					smtpMock.EXPECT().Rcpt(emailFromStr).After(firstCall).Return(nil).Times(1)

					return smtpMock
				},
			},
			args: args{
				addrs: addrs,
			},
			want: make(map[string]error),
		},
		{
			name: "error",
			fields: fields{
				client: func() smtpclient.SMTPClient {
					smtpMock := mocksmtpclient.NewMockSMTPClient(ctrl)
					firstCall := smtpMock.EXPECT().Rcpt(emailToStr).Return(errorSimple).Times(1)
					smtpMock.EXPECT().Rcpt(emailFromStr).After(firstCall).Return(nil).Times(1)

					return smtpMock
				},
			},
			args: args{
				addrs: addrs,
			},
			want: map[string]error{emailToStr: errorSimple},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{
				client:    tt.fields.client(),
				tlsConfig: tt.fields.TLSConfig,
			}
			if got := s.RCPTs(tt.args.addrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RCPTs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendMail_SetClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		client smtpclient.SMTPClient
	}
	type args struct {
		client interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   SendMail
	}{
		{
			name: "success",
			fields: fields{
				client: simpleClient,
			},
			args: args{
				client: simpleClient,
			},
			want: &sendMail{client: simpleClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sendMail{}
			if s.SetClient(tt.fields.client); !reflect.DeepEqual(tt.want, s) {
				t.Errorf("SetClient() server = %v, want %v", s, tt.want)
			}
		})
	}
}
