# λ v2

Functional programming helpers for Go — **typed pipelines** built around `Option[T]`.

Import:

```go
import λ "github.com/4thel00z/lambda/v2"
```

## Design notes (v2)

- `Option[T]` holds a value + error.
- Go does **not** allow methods with their own type parameters, and does **not** allow “specialized methods” on `Option[ConcreteType]`.
  - Type-changing operations are generic **functions** like `λ.Map`, `λ.Then`, `λ.Try`.
  - Type-specific fluent pipelines are implemented via small wrapper types like `λ.Bytes`, `λ.Str`, `λ.Req`, `λ.Resp`, `λ.RSAKeys`.

## Quickstart

### Read all stdin

```go
package main

import (
	"os"
	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	content := λ.Slurp(os.Stdin).String().Must()
	_ = content
}
```

### Read a file and write to stdout

```go
package main

import (
	"os"
	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	λ.Open("lorem_ipsum.txt").Slurp().WriteToWriter(os.Stdout)
}
```

### Read JSON into a struct

```go
package main

import (
	"fmt"
	λ "github.com/4thel00z/lambda/v2"
)

type MagicSpell struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

func main() {
	spell := λ.FromJSON[MagicSpell](λ.Open("magic.json").Slurp()).Must()
	fmt.Println(spell.Name, spell.Power)
}
```

### Simple HTTP request

```go
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
```

### Conditionals (kept, typed)

```go
package main

import (
	"errors"
	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	manipulateError := λ.Return(λ.Err[int](errors.New("this error will be thrown")))
	input := λ.Wrap(0, errors.New("something is weird"))
	output := λ.If(λ.HasError[int], manipulateError).Else(λ.Identity[int]).Do(input)
	_ = output
}
```


