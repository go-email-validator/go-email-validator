package evsmtp

import (
	"bytes"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/stretchr/testify/require"
	"io"
	"net/smtp"
	"reflect"
	"strings"
	"testing"
)

const (
	smSetClient      = "SetClient"
	smClient         = "Client"
	smHello          = "Hello "
	smHelloLocalhost = "Hello localhost"
	smAuth           = "Auth"
	smMail           = "Mail "
	smRCPTs          = "Rcpts "
	smData           = "Data"
	smWrite          = "Write"
	smQuit           = "Quit"
	smClose          = "Close"
	smWriteMessage   = "Write message"
	smWCloseWriter   = "Close writer"
)

var (
	testUser           = "testUser"
	testPwd            = "testPwd"
	testHost           = "testHost"
	testAuth           = smtp.PlainAuth("", testUser, testPwd, testHost)
	testMsg            = "msg"
	mockWriterInstance = &mockWriter{}
)

func stringsJoin(strs []string) string {
	return strings.Join(strs, ",")
}

// TODO create mock by gomock
type sendMailWant struct {
	stage   string
	message string
	ret     interface{}
}

type mockSendMail struct {
	t    *testing.T
	i    int
	want []sendMailWant
}

func (s *mockSendMail) SetClient(i interface{}) {
	require.Equal(s.t, s.do(smSetClient), i)
}

func (s *mockSendMail) Client() interface{} {
	return s.do(smClient)
}

func (s *mockSendMail) Hello(localName string) error {
	return evtests.ToError(s.do(smHello + localName))
}

func (s *mockSendMail) Auth(a smtp.Auth) error {
	ret := s.do(smAuth).([]interface{})
	if !reflect.DeepEqual(a, ret[0]) {
		s.t.Errorf("Invalid auth, got %#v, want %#v", a, testAuth)
	}
	return evtests.ToError(ret[1])
}

func (s *mockSendMail) Mail(from string) error {
	return evtests.ToError(s.do(smMail + from))
}

func (s *mockSendMail) RCPTs(addr []string) map[string]error {
	err := s.do(smRCPTs + stringsJoin(addr))

	if err == nil {
		return nil
	}

	return map[string]error{
		addr[0]: evtests.ToError(err),
	}
}

func (s *mockSendMail) Data() (io.WriteCloser, error) {
	return &mockWriter{s: s, want: testMsg}, evtests.ToError(s.do(smData))
}

func (s *mockSendMail) Write(w io.WriteCloser, msg []byte) error {
	w.Write(msg)
	w.Close()

	return evtests.ToError(s.do(smWrite))
}

func (s *mockSendMail) Quit() error {
	return evtests.ToError(s.do(smQuit))
}

func (s *mockSendMail) Close() error {
	return evtests.ToError(s.do(smClose))
}

func (s *mockSendMail) do(cmd string) interface{} {
	if s.i >= len(s.want) {
		s.t.Fatalf("Invalid command %q", cmd)
	}

	if cmd != s.want[s.i].message {
		s.t.Fatalf("Invalid command, got %q, want %q", cmd, s.want[s.i].message)
	}
	s.i++

	return s.want[s.i-1].ret
}

type mockWriter struct {
	want string
	s    *mockSendMail
	buf  bytes.Buffer
}

func (w *mockWriter) Write(p []byte) (int, error) {
	if w.buf.Len() == 0 {
		w.s.do(smWriteMessage)
	}
	w.buf.Write(p)
	return len(p), nil
}

func (w *mockWriter) Close() error {
	require.Equal(w.s.t, w.buf.String(), w.want)
	w.s.do(smWCloseWriter)
	return nil
}

var defaultWantMap = map[string]sendMailWant{
	smSetClient: {
		message: smSetClient,
		ret:     simpleClient,
	},
	smHello: {
		message: smHelloLocalhost,
	},
	smAuth: {
		message: smAuth,
		ret:     []interface{}{nil, nil},
	},
	smMail: {
		message: smMail + emailFrom.String(),
		ret:     nil,
	},
	smRCPTs: {
		message: smRCPTs + randomAddress.String(),
		ret:     nil,
	},
}

var quitStageWant = sendMailWant{
	message: smQuit,
	ret:     nil,
}

var closeStageWant = sendMailWant{
	message: smClose,
	ret:     nil,
}

var wantSuccessList = []string{
	smSetClient,
	smHello,
	smAuth,
	smMail,
	smRCPTs, // only random email call rcpt
	smQuit,
}

func failWant(failStage *sendMailWant, withClose bool) []sendMailWant {
	wants := make([]sendMailWant, 0)
	for _, stage := range wantSuccessList {

		var want sendMailWant
		want, ok := defaultWantMap[stage]
		if !ok {
			want = sendMailWant{message: stage}
		}

		if failStage != nil && stage == failStage.stage {
			wants = append(wants, *failStage)
			break
		}
		wants = append(wants, want)
	}

	if withClose {
		wants = append(wants, closeStageWant)
	}

	return wants
}
