package v1

import (
	"encoding/json"
	"reflect"
)

type JSONSetter interface {
	Set(string) (interface{}, error)
}

type JSONGetter interface {
	Get(string) (interface{}, error)
}

func (o Option) SetJSON(setter JSONSetter) Option {
	return Wrap(setter.Set(o.UnwrapString()))
}

func (o Option) GetJSON(getter JSONGetter) Option {
	return Wrap(getter.Get(o.UnwrapString()))
}

func (o Option) JSON(val interface{}) Option {
	var (
		err error
	)

	b := o.UnwrapBytes()

	if reflect.ValueOf(val).Kind() == reflect.Ptr {
		err = json.Unmarshal(b, val)
	} else {
		err = json.Unmarshal(b, &val)
	}

	return Option{
		value: val,
		err:   err,
	}
}

func (o Option) ToJSON() Option {
	var (
		b   []byte
		err error
	)

	val := o.Value()

	if reflect.ValueOf(val).Kind() == reflect.Ptr {
		b, err = json.Marshal(val)
	} else {
		b, err = json.Marshal(&val)
	}
	return Option{
		value: b,
		err:   err,
	}
}
