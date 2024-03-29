package v1

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Option struct {
	value interface{}
	err   error
}

type Producer func(o Option) interface{}
type Transformer func(o Option) Option

type ErrorHandler func(err error) error

func Return(option Option) Transformer {
	return func(o Option) Option {
		return option
	}
}

func Wrap(i interface{}, err error) Option {
	return Option{
		value: i,
		err:   err,
	}
}

func WrapValue(i interface{}) Option {
	return Option{
		value: i,
	}
}

func WrapError(err error) Option {
	return Option{
		err: err,
	}
}
func (o Option) Error() error {
	return o.err
}

func (o Option) Value() interface{} {
	return o.value
}

func (o Option) Or(i interface{}) interface{} {
	if o.err != nil {
		return i
	}
	return nil
}

func (o Option) Read() Option {
	if o.err != nil {
		return Read(o.value.(io.Reader))
	}
	return o
}
func (o Option) Slurp() Option {
	if o.err != nil {
		return o
	}
	switch o.value.(type) {
	case io.ReadCloser:
		return Slurp(o.value.(io.ReadCloser))
	case io.Reader:
		return Read(o.value.(io.Reader))
	case *http.Response:
		return Slurp(o.value.(*http.Response).Body)
	case http.Response:
		return Slurp(o.value.(http.Response).Body)
	}

	return WrapError(fmt.Errorf("couldn't slurp %#v", o))
}

func (o Option) ToString() Option {
	if o.err != nil {
		return o
	}

	return Option{
		value: string(o.value.([]byte)),
		err:   o.err,
	}
}

func (o Option) Unwrap() interface{} {
	if o.err != nil {
		log.Fatal(o.err)
	}
	return o.value
}

func (o Option) UnwrapBytes() []byte {
	if o.err != nil {
		log.Fatal(o.err)
	}
	return o.value.([]byte)
}

func (o Option) UnwrapBytesReader() io.Reader {
	return bytes.NewReader(o.UnwrapBytes())
}

func (o Option) UnwrapStringReader() io.Reader {
	return strings.NewReader(o.UnwrapString())
}

func (o Option) UnwrapChecksum() string {
	if o.err != nil {
		log.Fatal(o.err)
	}
	return fmt.Sprintf("%x", o.value.([sha256.Size]byte))
}

func (o Option) Unwrap224Checksum() string {
	if o.err != nil {
		log.Fatal(o.err)
	}
	return fmt.Sprintf("%x", o.value.([sha256.Size224]byte))

}

func (o Option) UnwrapString() string {
	if o.err != nil {
		log.Fatal(o.err)
	}

	if s, ok := o.value.(string); ok {
		return s
	}

	return string(o.UnwrapBytes())
}

func (o Option) Catch(e ErrorHandler) Option {
	if o.err != nil {
		return Option{
			value: o.value,
			err:   e(o.err),
		}
	}
	return o
}

func (o Option) Close() Option {
	if o.err == nil {
		o.err = o.value.(io.Closer).Close()
	}
	return o
}

func (o Option) WriteFromReader(r io.Reader) Option {
	toWrite, err := io.ReadAll(r)
	if err != nil {
		return Option{
			value: o.value,
			err:   err,
		}
	}
	w, ok := o.value.(io.Writer)
	if !ok {
		return Option{
			value: o.value,
			err:   errors.New("option is not a writer"),
		}
	}
	_, err = w.Write(toWrite)

	return Option{
		value: o.value,
		err:   err,
	}
}
func (o Option) WriteFromBytes(toWrite []byte) Option {
	w, ok := o.value.(io.Writer)
	if !ok {
		return Option{
			value: o.value,
			err:   errors.New("option is not a writer"),
		}
	}
	_, err := w.Write(toWrite)

	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) WriteFromString(toWrite string) Option {
	w, ok := o.value.(io.StringWriter)
	if !ok {
		return o.WriteFromBytes([]byte(toWrite))
	}
	_, err := w.WriteString(toWrite)

	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) WriteTo(w io.Writer) Option {
	_, err := w.Write(o.UnwrapBytes())
	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) CopyToWriter(w io.Writer) Option {
	if o.err != nil {
		return o
	}
	switch o.value.(type) {

	case io.Reader:
		_, err := io.Copy(w, o.value.(io.Reader))
		return Wrap(o.value, err)

	case *http.Response:
		_, err := io.Copy(w, o.value.(*http.Response).Body)
		return Wrap(o.value, err)
	case []byte:
		_, err := io.Copy(w, bytes.NewReader(o.value.([]byte)))
		return Wrap(o.value, err)

	}
	return Wrap(o.value, fmt.Errorf("could not handle value of %s type yet", reflect.TypeOf(o.value).Name()))
}

func (o Option) WriteStringTo(w io.StringWriter) Option {
	_, err := w.WriteString(o.UnwrapString())
	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) Map(transformer Transformer) Option {
	return transformer(o)
}

func (o Option) Apply(producer Producer) interface{} {
	return producer(o)
}
