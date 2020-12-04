package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	λ.Open("lorem_ipsum.txt").Slurp().WriteStringTo(os.Stdout)
}
