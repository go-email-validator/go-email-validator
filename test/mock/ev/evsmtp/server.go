package mockevsmtp

import (
	"bufio"
	"bytes"
	"net"
	"net/textproto"
	"sync"
	"testing"
	"time"
)

// Separator separate mock message of Server
const Separator = "\r\n"

// Server to testing SMTP
// Partial copy of TestSendMail  from smtp.TestSendMail
func Server(t *testing.T, server []string, timeout time.Duration) (string, chan string) {
	var cmdbuf bytes.Buffer
	bcmdbuf := bufio.NewWriter(&cmdbuf)
	l, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatalf("Unable to create listener: %v", err)
	}

	var done = make(chan string)
	closedMu := &sync.Mutex{}
	closed := false
	closeServer := func() {
		closedMu.Lock()
		if !closed {
			closed = true
			bcmdbuf.Flush()
			done <- cmdbuf.String()
			close(done)
			l.Close()
		}
		closedMu.Unlock()
	}
	go func(data []string) {
		defer closeServer()

		if len(data) == 0 {
			return
		}

		conn, err := l.Accept()
		if err != nil {
			t.Errorf("Accept error: %v", err)
			return
		}
		defer conn.Close()

		tc := textproto.NewConn(conn)
		for i := 0; i < len(data) && data[i] != ""; i++ {
			tc.PrintfLine(data[i])
			for len(data[i]) >= 4 && data[i][3] == '-' {
				i++
				tc.PrintfLine(data[i])
			}
			if data[i] == "221 Goodbye" {
				return
			}
			read := false
			for !read || data[i] == "354 Go ahead" {
				msg, err := tc.ReadLine()
				bcmdbuf.Write([]byte(msg + Separator))
				read = true
				if err != nil {
					t.Errorf("Read error: %v", err)
					return
				}
				if data[i] == "354 Go ahead" && msg == "." {
					break
				}
			}
		}
	}(server)

	go func() {
		time.Sleep(timeout)
		closeServer()
	}()

	return l.Addr().String(), done
}
