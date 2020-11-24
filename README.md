# λ

![lambda-tests](https://github.com/4thel00z/lambda/workflows/Test/badge.svg)
![lambda-logo](https://github.com/4thel00z/lambda/raw/assets/logo.svg)

λ is a functional programming framework for go, which adds support for an alternative error handling workflow using options.

## Example usage

This demonstrates how you can use lambda in your code: 

```go
package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	fmt.Print(string(λ.Slurp(os.Stdin).UnwrapBytes()))
}

```

## Todo

* TBD

## License

This project is licensed under the GPL-3 license.
