package evsmtp

import (
	"errors"
	"fmt"
)

const (
	ErrorHelloAfter = "smtp: Hello called after other methods"
	ErrorCrLR       = "smtp: A line must not contain CR or LF"
)

type Error interface {
	error
	Stage() SendMailStage
	Unwrap() error
}

type ASMTPError struct {
	stage SendMailStage
	err   error
}

func (a ASMTPError) Stage() SendMailStage {
	return a.stage
}

func (a ASMTPError) Unwrap() error {
	return a.err
}

func (a ASMTPError) Error() string {
	return fmt.Sprintf("%v happened on stage \"%v\"", errors.Unwrap(a).Error(), a.Stage())
}

func NewError(stage SendMailStage, err error) Error {
	return DefaultError{ASMTPError{stage, err}}
}

type DefaultError struct {
	ASMTPError
}
