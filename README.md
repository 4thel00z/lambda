# λ

![lambda-tests](https://github.com/4thel00z/lambda/workflows/Test/badge.svg)
![lambda-logo](https://media.githubusercontent.com/media/4thel00z/lambda/master/logo.png)

λ is a functional programming framework for go, which adds support for an alternative error handling workflow using options.

## Example usage


### Read all lines from stdin

```go
package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	content := λ.Slurp(os.Stdin).UnwrapString()
	// do things with content...
}


```

### Read a file and pipe it to stdout

```go
package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	λ.Open("lorem_ipsum.txt").Slurp().WriteString(os.Stdout)
}
```

### Read a JSON file into a struct

```go
package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"strings"
)

type MagicSpell struct {
	Name        string  `json:"name"`
	AttackPower float64 `json:"attack_power"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
}

func main() {
	var (
		m MagicSpell
	)
	λ.Open("magic.json").Slurp().JSON(&m).Catch(λ.Die)
}

```

### Functional conditionals

You never need to check an error with an if clause again. Instead you can define the flow as functional chain,
start point is always `λ.If`.
You even can reuse the same chain, it doesn't contain data. You pass the data via `Conditional.Do`.

```go
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

```

### Make Rest calls

```go
package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	λ.Get("https://ransomware.host").Do().Slurp().WriteString(os.Stdout)
}
```

## Todo

* Make Option more flexible an pretty

## License

This project is licensed under the GPL-3 license.
