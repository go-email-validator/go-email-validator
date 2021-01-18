package ev

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"hash"
	"net/http"
)

const (
	// GravatarValidatorName is name for validation by https://www.gravatar.com/
	GravatarValidatorName ValidatorName = "Gravatar"
	// GravatarURL is url for gravatar validation
	GravatarURL string = "https://www.gravatar.com/avatar/%x?d=404"
)

// GravatarError is GravatarValidatorName error
type GravatarError struct{}

func (GravatarError) Error() string {
	return "GravatarError"
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

	g.h.Reset()
	g.h.Write([]byte(input.Email().String()))
	resp, err := http.Head(fmt.Sprintf(GravatarURL, g.h.Sum(nil)))
	if err != nil || resp.StatusCode != 200 {
		return gravatarGetError(GravatarError{})
	}
	defer resp.Body.Close()

	return NewValidResult(GravatarValidatorName)
}

func gravatarGetError(err error) ValidationResult {
	return NewResult(false, utils.Errs(err), nil, GravatarValidatorName)
}
