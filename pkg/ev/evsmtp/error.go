package evsmtp

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"net/textproto"
	"reflect"
)

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

type Error interface {
	error
	Stage() SendMailStage
	Unwrap() error
}

type AliasError error

type ASMTPError struct {
	stage SendMailStage
	err   error
}

func (a *ASMTPError) Stage() SendMailStage {
	return a.stage
}

func (a *ASMTPError) Unwrap() error {
	return a.err
}

func (a *ASMTPError) Error() string {
	return fmt.Sprintf("%v happened on stage \"%v\"", errors.Unwrap(a).Error(), a.Stage())
}

func (a *ASMTPError) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.EncodeMulti(a.stage, a.err)
}

func (a *ASMTPError) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.DecodeMulti(&a.stage, &a.err)
}

func NewError(stage SendMailStage, err error) Error {
	return &DefaultError{ASMTPError{stage, err}}
}

type DefaultError struct {
	ASMTPError
}

func ConvertEVSMTPErrorsToErrors(Errs []AliasError) (errs []error) {
	errs = make([]error, len(Errs))
	for i, Err := range Errs {
		errs[i] = Err
	}

	return
}

func ConvertErrorsToEVSMTPErrors(errs []error) (Errs []AliasError) {
	Errs = make([]AliasError, len(errs))
	for i, err := range errs {
		Errs[i] = err
	}

	return
}
