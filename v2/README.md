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

## Parallelism (multi-threading)

Parallel helpers for slices using `errgroup` with a configurable concurrency limit.

### Parallel map (ordered)

```go
package main

import (
	"context"
	"fmt"
	"strings"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	ctx := context.Background()
	words := []string{"a", "bb", "ccc"}

	out := λ.ParMap(ctx, words, func(s string) string {
		return strings.ToUpper(s)
	}, λ.WithConcurrency(8)).Must()

	fmt.Println(out) // [A BB CCC]
}
```

### Parallel try-map (fail-fast)

```go
package main

import (
	"context"
	"errors"
	"strconv"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	ctx := context.Background()
	in := []string{"1", "2", "nope", "4"}

	parse := λ.TryFn[string, int](func(s string) (int, error) {
		if s == "nope" {
			return 0, errors.New("bad input")
		}
		return strconv.Atoi(s)
	})

	_, _ = λ.ParTry(ctx, in, parse, λ.WithConcurrency(4)).Get()
}
```

### Parallel map over a channel (unordered)

```go
package main

import (
	"context"
	"fmt"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	ctx := context.Background()
	in, inErrc := λ.RangeN(ctx, 5)

	out, errc := λ.ParMapChan(ctx, in, func(v int) int { return v * 2 }, λ.WithConcurrency(8))

	for v := range out {
		fmt.Println(v)
	}
	_ = λ.JoinErr(inErrc, errc) // nil on success
}
```

## Channel utilities

Exporters like `RangeN` and transforms like `Take` let you build channel flows without boilerplate goroutines:

```go
package main

import (
	"context"
	"fmt"

	λ "github.com/4thel00z/lambda/v2"
)

func main() {
	ctx := context.Background()

	src, srcErrc := λ.RangeN(ctx, 10)
	first3, first3Errc := λ.Take(ctx, src, 3)

	fmt.Println(λ.Collect(ctx, first3).Must()) // [0 1 2]
	_ = <-srcErrc
	_ = <-first3Errc
}
```


