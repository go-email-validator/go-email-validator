package evsmtp

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"net/textproto"
	"reflect"
)

// Constants of smtpclient.SMTPClient errors
const (
	ErrorHelloAfter = "smtp: Hello called after other methods"
	ErrorCrLR       = "smtp: A line must not contain CR or LF"
)

func init() {
	msgpack.RegisterExt(1, new(DefaultError))
	msgpack.RegisterExt(2, new(ASMTPError))
	msgpack.RegisterExt(3, new(textproto.Error))

	msgpack.Register(errors.New(""), func(e *msgpack.Encoder, v reflect.Value) error {
		if v.IsNil() {
			return e.EncodeNil()
		}
		return e.EncodeString(v.Interface().(error).Error())
	}, nil)
}

// Error is interface of Checker errors
type Error interface {
	error
	Stage() SendMailStage
	Unwrap() error
}

// AliasError is alias to fix msgpack
type AliasError error

// ASMTPError isa abstract struct for Checker errors
type ASMTPError struct {
	stage SendMailStage
	err   error
}

// Stage returns stage of error
func (a *ASMTPError) Stage() SendMailStage {
	return a.stage
}

func (a *ASMTPError) Unwrap() error {
	return a.err
}

func (a *ASMTPError) Error() string {
	return fmt.Sprintf("%v happened on stage \"%v\"", errors.Unwrap(a).Error(), a.Stage())
}

// EncodeMsgpack implements encoder for msgpack
func (a *ASMTPError) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(a.stage, a.err)
}

// DecodeMsgpack implements decoder for msgpack
func (a *ASMTPError) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&a.stage, &a.err)
}

// NewError is constructor for DefaultError
func NewError(stage SendMailStage, err error) Error {
	return &DefaultError{ASMTPError{stage, err}}
}

// DefaultError is default error
type DefaultError struct {
	ASMTPError
}

// Convert []AliasError to []error
func _(Errs []AliasError) (errs []error) {
	errs = make([]error, len(Errs))
	for i, Err := range Errs {
		errs[i] = Err
	}

	return
}

// ErrorsToEVSMTPErrors converts []error to []AliasError
// It is used like fix of msgpack problems https://github.com/vmihailenco/msgpack/issues/294
func ErrorsToEVSMTPErrors(errs []error) (Errs []AliasError) {
	Errs = make([]AliasError, len(errs))
	for i, err := range errs {
		Errs[i] = err
	}

	return
}
