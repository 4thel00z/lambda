package v1

import (
	"encoding/json"
	"io"
)

type Option struct {
	value interface{}
	err   error
}

type Producer func(o Option) interface{}
type ErrorHandler func(err error)

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
	return Slurp(o.value.(io.ReadCloser))

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
		panic(o.err)
	}
	return o.value
}

func (o Option) UnwrapBytes() []byte {
	if o.err != nil {
		panic(o.err)
	}
	return o.value.([]byte)
}

func (o Option) UnwrapString() string {
	if o.err != nil {
		panic(o.err)
	}
	return string(o.UnwrapBytes())
}

func (o Option) Catch(e ErrorHandler) interface{} {
	if o.err != nil {
		e(o.err)
		return nil
	}
	return o.value
}

func (o Option) Close() Option {
	if o.err == nil {
		o.err = o.value.(io.Closer).Close()
	}
	return o
}

func (o Option) Write(w io.Writer) Option {
	_, err := w.Write(o.UnwrapBytes())
	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) WriteString(w io.StringWriter) Option {
	_, err := w.WriteString(o.UnwrapString())
	return Option{
		value: o.value,
		err:   err,
	}
}

func (o Option) JSON(i interface{}) Option {
	b := o.UnwrapBytes()
	err := json.Unmarshal(b, &i)
	return Option{
		value: i,
		err:   err,
	}
}
