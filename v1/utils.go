package v1

import (
	"errors"
	"log"
)

// Some convenience ErrorHandler

func Cry(o Option) Option { log.Fatalln(o.Error()); return o }
func Die(err error) error { log.Fatalln(err); return nil }
func Ignore(error) error  { return nil }

var (
	Error = errors.New
)
