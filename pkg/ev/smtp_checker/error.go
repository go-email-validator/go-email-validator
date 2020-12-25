package smtp_checker

import "fmt"

const (
	SMTPErrorHelloAfter = "smtp_checker: Hello called after other methods"
	SMTPErrorCrLR       = "smtp_checker: A line must not contain CR or LF"
)

type SMTPError interface {
	error
	Stage() SendMailStage
	Err() error
}

type ASMTPError struct {
	stage SendMailStage
	err   error
}

func (a ASMTPError) Stage() SendMailStage {
	return a.stage
}

func (a ASMTPError) Err() error {
	return a.err
}

func (a ASMTPError) Error() string {
	return fmt.Sprintf("%v happend on stage \"%v\"", a.Err().Error(), a.Stage())
}

func NewSmtpError(stage SendMailStage, err error) SMTPError {
	return DefaultSmtpError{ASMTPError{stage, err}}
}

type DefaultSmtpError struct {
	ASMTPError
}

type SMTPErrorNested interface {
	SMTPError
	GetNested() SMTPError
}

type ASMTPErrorNested struct {
	n SMTPError
}

func (a ASMTPErrorNested) GetNested() SMTPError {
	return a.n
}

func (a ASMTPErrorNested) Error() string {
	return a.n.Error()
}
