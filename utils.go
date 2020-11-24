package lambda

import "log"

func Die(err error) { log.Fatalln(err) }
func Ignore(error)  {}
