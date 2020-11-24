package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	fmt.Print(string(λ.Slurp(os.Stdin).UnwrapBytes()))
}
