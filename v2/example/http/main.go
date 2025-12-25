package main

import (
	"context"
	"os"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	λ.Get("https://example.com").
		Do(context.Background()).
		Slurp().
		WriteToWriter(os.Stdout)
}


