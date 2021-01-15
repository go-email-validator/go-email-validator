package evsmtp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"net"
	"net/textproto"
	"net/url"
	"reflect"
)

// Constants of smtpclient.SMTPClient errors
const (
	ErrorHelloAfter = "smtp: Hello called after other methods"
	ErrorCrLR       = "smtp: A line must not contain CR or LF"
)

// Used 0 because of https://github.com/msgpack/msgpack/blob/master/spec.md#extension-types
var registerExtId int8 = 0

// ExtId returns register extent id, used for msgpack.RegisterExt
func ExtId() int8 {
	registerExtId++
	return registerExtId - 1
}

// SetExtId sets register extent id, used for msgpack.RegisterExt
func SetExtId(rExt int8) {
	registerExtId = rExt
}

// Import different error types from packages, used in smtp.Client
func init() {
	msgpack.RegisterExt(ExtId(), new(DefaultError))
	msgpack.RegisterExt(ExtId(), new(ASMTPError))

	msgpack.RegisterExt(ExtId(), new(textproto.Error))
	msgpack.RegisterExt(ExtId(), new(textproto.ProtocolError))

	msgpack.RegisterExt(ExtId(), new(net.AddrError))
	msgpack.RegisterExt(ExtId(), new(net.DNSConfigError))
	msgpack.RegisterExt(ExtId(), new(net.DNSError))
	msgpack.RegisterExt(ExtId(), new(net.InvalidAddrError))
	msgpack.RegisterExt(ExtId(), new(net.OpError))
	msgpack.RegisterExt(ExtId(), new(net.ParseError))
	msgpack.RegisterExt(ExtId(), new(net.UnknownNetworkError))

	msgpack.RegisterExt(ExtId(), new(url.Error))
	msgpack.RegisterExt(ExtId(), new(url.EscapeError))
	msgpack.RegisterExt(ExtId(), new(url.InvalidHostError))

	msgpack.RegisterExt(ExtId(), new(tls.RecordHeaderError))
	msgpack.RegisterExt(ExtId(), new(x509.CertificateInvalidError))
	msgpack.RegisterExt(ExtId(), new(x509.HostnameError))
	msgpack.RegisterExt(ExtId(), new(x509.UnknownAuthorityError))
	msgpack.RegisterExt(ExtId(), new(x509.SystemRootsError))
	msgpack.RegisterExt(ExtId(), new(x509.InsecureAlgorithmError))
	msgpack.RegisterExt(ExtId(), new(x509.ConstraintViolationError))

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
