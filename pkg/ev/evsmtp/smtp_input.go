package evsmtp

import (
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"time"
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
	Timeout() time.Duration
	Port() int
}

// NewInput instantiates Input
func NewInput(email evmail.Address, opts Options) Input {
	if opts == nil {
		opts = DefaultOptions()
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
	EmailFrom evmail.Address
	HelloName string
	Proxy     string
	Timeout   time.Duration
	Port      int
}

var defaultOptions = NewOptions(OptionsDTO{})

func DefaultOptions() Options {
	return defaultOptions
}

// NewOptions instantiates Options
func NewOptions(dto OptionsDTO) Options {
	return &options{
		emailFrom: dto.EmailFrom,
		helloName: dto.HelloName,
		proxy:     dto.Proxy,
		timeout:   dto.Timeout,
		port:      dto.Port,
	}
}

type options struct {
	emailFrom evmail.Address
	helloName string
	proxy     string
	timeout   time.Duration
	port      int
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
func (i *options) Timeout() time.Duration {
	return i.timeout
}
func (i *options) Port() int {
	return i.port
}
