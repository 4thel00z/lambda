package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	fmt.Print(λ.Open("lorem_ipsum.txt").Slurp().ToString().UnwrapString())
}
