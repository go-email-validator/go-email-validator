package ev

import (
	"bitbucket.org/maranqz/email-validator/pkg/ev/ev_email"
	"context"
	"fmt"
	"github.com/smancke/mailck"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const SMTPValidatorName = "SMTPValidatorInterface"

type SMTPValidatorInterface interface {
	DepValidatorInterface
}

type SMTPValidator struct {
	port      uint16
	fromEmail ev_email.EmailAddressInterface
	helloName string
	ADepValidator
}

func NewSMTPValidator(fromEmail ev_email.EmailAddressInterface, port *uint16) *SMTPValidator {
	if fromEmail == nil {
		fromEmail = ev_email.NewEmailAddress("user", "example.org")
	}
	if port == nil {
		port = new(uint16)
		*port = 25
	}

	return &SMTPValidator{
		*port,
		fromEmail,
		"localhost",
		ADepValidator{},
	}
}

func (a *SMTPValidator) GetDeps() []string {
	return []string{SyntaxValidatorName, MXValidatorName}
}

func (a *SMTPValidator) Validate(email ev_email.EmailAddressInterface) ValidationResultInterface {
	var err error = nil
	var result mailck.Result
	syntaxResult := (*a.results)[0].(SyntaxValidatorResultInterface)
	mxResult := (*a.results)[1].(MXValidationResultInterface)

	if syntaxResult.IsValid() && mxResult.IsValid() {
		result, err = checkMailbox(mxResult.MX(), email.String(), a.fromEmail.String())
	}

	return NewValidatorResult(err == nil && result.IsValid(), nil, nil)
}

var noContext = context.Background()
var defaultDialer = newDialer()

func newDialer() net.Dialer {
	var defaultDialer = net.Dialer{}
	defaultDialer.Timeout = 2 * time.Second
	return defaultDialer
}

func checkMailbox(mxList MXs, checkEmail, fromEmail string) (result mailck.Result, err error) {
	if len(mxList) == 0 {
		return mailck.InvalidDomain, nil
	}
	return checkMailboxMailck(noContext, fromEmail, checkEmail, mxList, 25)
}

type checkRv struct {
	res mailck.Result
	err error
}

// get from https://github.com/FGRibreau/mailchecker/blob/master/platform/go/mail_checker.go
func checkMailboxMailck(ctx context.Context, fromEmail, checkEmail string, mxList []*net.MX, port int) (result mailck.Result, err error) {
	// try to connect to one mx
	var c *smtp.Client

	for _, mx := range mxList {
		var conn net.Conn
		conn, err = defaultDialer.DialContext(ctx, "tcp", fmt.Sprintf("%v:%v", mx.Host, port))
		if t, ok := err.(*net.OpError); ok {
			if t.Timeout() {
				return mailck.TimeoutError, err
			}
			return mailck.NetworkError, err
		} else if err != nil {
			return mailck.MailserverError, err
		}
		c, err = smtp.NewClient(conn, mx.Host)
		if err == nil {
			break
		}
	}
	if err != nil {
		return mailck.MailserverError, err
	}
	if c == nil {
		// just to get very sure, that we have a connection
		// this code line should never be reached!
		return mailck.MailserverError, fmt.Errorf("can't obtain connection for %v", checkEmail)
	}

	resChan := make(chan checkRv, 1)

	go func() {
		defer c.Close()
		defer c.Quit() // defer ist LIFO
		// HELO
		err = c.Hello(hostname(fromEmail))
		if err != nil {
			resChan <- checkRv{mailck.MailserverError, err}
			return
		}

		// MAIL FROM
		err = c.Mail(fromEmail)
		if err != nil {
			resChan <- checkRv{mailck.MailserverError, err}
			return
		}

		// RCPT TO
		id, err := c.Text.Cmd("RCPT TO:<%s>", checkEmail)
		if err != nil {
			resChan <- checkRv{mailck.MailserverError, err}
			return
		}
		c.Text.StartResponse(id)
		code, message, err := c.Text.ReadResponse(25)
		c.Text.EndResponse(id)
		_ = message
		if code == 550 {
			resChan <- checkRv{mailck.MailboxUnavailable, nil}
			return
		}

		if err != nil {
			resChan <- checkRv{mailck.MailserverError, err}
			return
		}

		resChan <- checkRv{mailck.Valid, nil}

	}()

	select {
	case <-ctx.Done():
		return mailck.TimeoutError, ctx.Err()
	case q := <-resChan:
		return q.res, q.err
	}
}

func hostname(mail string) string {
	return mail[strings.Index(mail, "@")+1:]
}
