package main

import (
	"errors"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	manipulateError := λ.Return(λ.Err[int](errors.New("this error will be thrown")))
	input := λ.Wrap(0, errors.New("something is weird"))
	_ = λ.If(λ.HasError[int], manipulateError).Else(λ.Identity[int]).Do(input)
}


