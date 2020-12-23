package ev

import (
	"crypto/md5"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/ev_email"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"hash"
	"net/http"
)

const (
	GravatarValidatorName ValidatorName = "Gravatar"
	GravatarUrl                         = "https://www.gravatar.com/avatar/%x?d=404"
)

type GravatarError struct {
	utils.Err
}

func NewGravatarValidator() Validator {
	return gravatarValidator{h: md5.New()}
}

type gravatarValidator struct {
	h hash.Hash
	AValidatorWithoutDeps
}

func (_ gravatarValidator) GetDeps() []ValidatorName {
	return []ValidatorName{SyntaxValidatorName}
}

func (w gravatarValidator) Validate(email ev_email.EmailAddress, results ...ValidationResult) ValidationResult {
	syntaxResult := results[0].(SyntaxValidatorResultInterface)
	if !syntaxResult.IsValid() {
		return gravatarGetError()
	}

	w.h.Reset()
	w.h.Write([]byte(email.String()))
	resp, err := http.Head(fmt.Sprintf(GravatarUrl, w.h.Sum(nil)))
	if err != nil || resp.StatusCode != 200 {
		return gravatarGetError()
	}
	defer resp.Body.Close()

	return NewValidValidatorResult(GravatarValidatorName)
}

func gravatarGetError() ValidationResult {
	return NewValidatorResult(false, utils.Errs(GravatarError{}), nil, GravatarValidatorName)
}
