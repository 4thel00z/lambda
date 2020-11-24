package lambda

import (
	"io"
	"io/ioutil"
)

func Slurp(r io.Reader) Option {
	content, err := ioutil.ReadAll(r)
	return Option{
		value: content,
		err:   err,
	}
}
