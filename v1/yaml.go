package v1

import (
	"gopkg.in/yaml.v3"
	"reflect"
)

type YAMLSetter interface {
	Set(string) (interface{}, error)
}

type YAMLGetter interface {
	Get(string) (interface{}, error)
}

func (o Option) SetYAML(setter YAMLSetter) Option {
	return Wrap(setter.Set(o.UnwrapString()))
}

func (o Option) GetYAML(getter YAMLGetter) Option {
	return Wrap(getter.Get(o.UnwrapString()))
}

func (o Option) YAML(val interface{}) Option {
	var (
		err error
	)

	b := o.UnwrapBytes()

	if reflect.ValueOf(val).Kind() == reflect.Ptr {
		err = yaml.Unmarshal(b, val)
	} else {
		err = yaml.Unmarshal(b, &val)
	}

	return Option{
		value: val,
		err:   err,
	}
}

func (o Option) ToYAML() Option {
	var (
		b   []byte
		err error
	)

	val := o.Value()

	if reflect.ValueOf(val).Kind() == reflect.Ptr {
		b, err = yaml.Marshal(val)
	} else {
		b, err = yaml.Marshal(&val)
	}
	return Option{
		value: b,
		err:   err,
	}
}
