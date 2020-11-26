package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	λ.Get("https://ransomware.host").Do().Slurp().WriteString(os.Stdout)
}
