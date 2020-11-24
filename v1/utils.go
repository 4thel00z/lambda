package v1

import "log"

// Some convenience ErrorHandler

func Die(err error) error { log.Fatalln(err); return nil }
func Ignore(error) error  { return nil }
