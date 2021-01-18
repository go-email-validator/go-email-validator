package evsmtp

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"time"
)

const (
	DefaultTimeoutConnection = 2 * time.Second
	DefaultTimeoutResponse   = 2 * time.Second
)

// Input describes data for Checker
type Input interface {
	Email() evmail.Address
	Options
}

// Options describes smtp options
type Options interface {
	EmailFrom() evmail.Address
	HelloName() string
	Proxy() string
	TimeoutConnection() time.Duration
	TimeoutResponse() time.Duration
	Port() int
}

// NewInput instantiates Input
func NewInput(email evmail.Address, opts Options) Input {
	if opts == nil {
		opts = EmptyOptions()
	}

	return &input{
		email:   email,
		Options: opts,
	}
}

type input struct {
	email evmail.Address
	Options
}

func (i *input) Email() evmail.Address {
	return i.email
}

// OptionsDTO is dto for NewOptions
type OptionsDTO struct {
	EmailFrom   evmail.Address
	HelloName   string
	Proxy       string
	TimeoutCon  time.Duration
	TimeoutResp time.Duration
	Port        int
}

var defaultOptions = NewOptions(OptionsDTO{
	TimeoutCon:  DefaultTimeoutConnection,
	TimeoutResp: DefaultTimeoutResponse,
})

// DefaultOptions returns options with default values
func DefaultOptions() Options {
	return defaultOptions
}

var emptyOptions = NewOptions(OptionsDTO{})

// EmptyOptions returns empty options to avoid rewriting of default values
func EmptyOptions() Options {
	return emptyOptions
}

// NewOptions instantiates Options
func NewOptions(dto OptionsDTO) Options {
	return &options{
		emailFrom:   dto.EmailFrom,
		helloName:   dto.HelloName,
		proxy:       dto.Proxy,
		timeoutCon:  dto.TimeoutCon,
		timeoutResp: dto.TimeoutResp,
		port:        dto.Port,
	}
}

type options struct {
	emailFrom   evmail.Address
	helloName   string
	proxy       string
	timeoutCon  time.Duration
	timeoutResp time.Duration
	port        int
}

func (i *options) EmailFrom() evmail.Address {
	return i.emailFrom
}
func (i *options) HelloName() string {
	return i.helloName
}
func (i *options) Proxy() string {
	return i.proxy
}
func (i *options) TimeoutConnection() time.Duration {
	return i.timeoutCon
}
func (i *options) TimeoutResponse() time.Duration {
	return i.timeoutResp
}
func (i *options) Port() int {
	return i.port
}
