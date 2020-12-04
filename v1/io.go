package v1

import (
	"io"
	"io/ioutil"
	"os"
)

func Slurp(r io.ReadCloser) Option {
	content, err := ioutil.ReadAll(r)
	defer func() {
		if err != nil {
			r.Close()
		}
	}()

	return Option{
		value: content,
		err:   err,
	}
}

func Read(r io.Reader) Option {
	content, err := ioutil.ReadAll(r)
	return Option{
		value: content,
		err:   err,
	}
}

func Open(path string) Option {
	f, err := os.Open(path)
	return Option{
		value: f,
		err:   err,
	}
}

func Create(path string) Option {
	f, err := os.Create(path)
	return Option{
		value: f,
		err:   err,
	}
}
