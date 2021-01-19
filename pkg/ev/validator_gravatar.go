package ev

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"hash"
	"net/http"
	"time"
)

const (
	// GravatarValidatorName is name for validation by https://www.gravatar.com/
	GravatarValidatorName ValidatorName = "Gravatar"
	// GravatarURL is url for gravatar validation
	GravatarURL string = "https://www.gravatar.com/avatar/%x?d=404"
	// DefaultGravatarTimeout is default timeout for gravatar connection
	DefaultGravatarTimeout = 5 * time.Second
)

// GravatarErr is text for GravatarError.Error
const GravatarErr = "GravatarError"

// GravatarError is GravatarValidatorName error
type GravatarError struct{}

func (GravatarError) Error() string {
	return GravatarErr
}

// NewGravatarValidator instantiates GravatarValidatorName validator with GravatarURL for validation
func NewGravatarValidator() Validator {
	return NewGravatarValidatorWithURL(GravatarURL)
}

// NewGravatarValidatorWithURL instantiates GravatarValidatorName validator
func NewGravatarValidatorWithURL(gravatarURL string) Validator {
	return gravatarValidator{
		h:   md5.New(), //nolint:gosec
		url: gravatarURL,
	}
}

type gravatarValidator struct {
	AValidatorWithoutDeps
	h   hash.Hash
	url string
}

func (g gravatarValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName}
}

func (g gravatarValidator) Validate(input Input, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResult)
	if !syntaxResult.IsValid() {
		return gravatarGetError(NewDepsError())
	}

	var opts = DefaultGravatarOptions()
	if optsInterface := input.Option(GravatarValidatorName); optsInterface != nil {
		opts = optsInterface.(GravatarOptions)
	}

	client := &http.Client{Timeout: opts.Timeout()}

	g.h.Reset()
	g.h.Write([]byte(input.Email().String()))
	resp, err := client.Head(fmt.Sprintf(GravatarURL, g.h.Sum(nil)))
	if err != nil {
		return gravatarGetError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return gravatarGetError(GravatarError{})
	}

	return NewValidResult(GravatarValidatorName)
}

func gravatarGetError(err error) ValidationResult {
	return NewResult(false, utils.Errs(err), nil, GravatarValidatorName)
}

// GravatarOptions describes gravatar options
type GravatarOptions interface {
	Timeout() time.Duration
}

// GravatarOptionsDTO is dto for NewGravatarOptions
type GravatarOptionsDTO struct {
	Timeout time.Duration
}

var defaultOptions = NewGravatarOptions(GravatarOptionsDTO{
	Timeout: DefaultGravatarTimeout,
})

// DefaultOptions returns options with default values
func DefaultGravatarOptions() GravatarOptions {
	return defaultOptions
}

// NewGravatarOptions instantiates GravatarOptions
func NewGravatarOptions(dto GravatarOptionsDTO) GravatarOptions {
	return &gravatarOptions{
		timeout: dto.Timeout,
	}
}

type gravatarOptions struct {
	timeout time.Duration
}

func (i *gravatarOptions) Timeout() time.Duration {
	return i.timeout
}
