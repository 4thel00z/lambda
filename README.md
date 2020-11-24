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

## Todo

* Make Option more flexible an pretty
* Add methods for handling conditionals 

## License

This project is licensed under the GPL-3 license.
