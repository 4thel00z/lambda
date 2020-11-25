package main

import (
	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	manipulateError := λ.Return(λ.Wrap(nil, λ.Error("this error will be thrown")))
	input := λ.Wrap(nil, λ.Error("something is weird"))
	output := λ.If(λ.HasError, manipulateError).Else(λ.Cry).Do(input)
	λ.If(λ.HasNoError, λ.Identity).Else(λ.Cry).Do(output)
}
