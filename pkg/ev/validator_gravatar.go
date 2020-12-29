package ev

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"hash"
	"net/http"
)

const (
	GravatarValidatorName ValidatorName = "Gravatar"
	GravatarURL           string        = "https://www.gravatar.com/avatar/%x?d=404"
)

type GravatarError struct {
	utils.Err
}

func NewGravatarValidator() Validator {
	return gravatarValidator{h: md5.New()} //nolint:gosec
}

type gravatarValidator struct {
	AValidatorWithoutDeps
	h hash.Hash
}

func (g gravatarValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName}
}

func (g gravatarValidator) Validate(email evmail.Address, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResult)
	if !syntaxResult.IsValid() {
		return gravatarGetError(DepsError{})
	}

	g.h.Reset()
	g.h.Write([]byte(email.String()))
	resp, err := http.Head(fmt.Sprintf(GravatarURL, g.h.Sum(nil)))
	if err != nil || resp.StatusCode != 200 {
		return gravatarGetError(GravatarError{})
	}
	defer resp.Body.Close()

	return NewValidValidatorResult(GravatarValidatorName)
}

func gravatarGetError(err error) ValidationResult {
	return NewValidatorResult(false, utils.Errs(err), nil, GravatarValidatorName)
}
