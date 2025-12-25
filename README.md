<div align="center">
  <img src="logo.png" width="160" alt="lambda logo" />

  <h1>λ</h1>

  <p><strong>Functional programming helpers for Go</strong> — Option-based error handling and fluent pipelines.</p>

  <p>
    <a href="https://github.com/4thel00z/lambda/actions/workflows/run_go_tests.yml">
      <img alt="CI" src="https://github.com/4thel00z/lambda/actions/workflows/run_go_tests.yml/badge.svg" />
    </a>
    <img alt="Go" src="https://img.shields.io/badge/go-1.23%2B-00ADD8.svg" />
    <a href="COPYING">
      <img alt="License: GPL-3.0" src="https://img.shields.io/badge/license-GPL--3.0-blue.svg" />
    </a>
  </p>
</div>

## Why

Go already has great error handling, but sometimes you want to express flows as **pipelines** instead of nested `if err != nil` blocks.
λ gives you an `Option` type (value + error) and a set of utilities that make those pipelines ergonomic.

## Requirements

- Go **1.23+** (this module tracks the `go` version in `go.mod`).

## Install


### v2 (default)

**λ v2 is a major rewrite** with generics + typed pipelines.

- Import:

```go
import λ "github.com/4thel00z/lambda/v2"
```

- This is a **hard cut**: **no compatibility layer** and **no migration path** from `v1`.

### v1 

Import the default versioned package:

```go
import λ "github.com/4thel00z/lambda/v1"
```

## Quickstart

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

### Parallelism (multi-threading)

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
	λ.Must(λ.JoinErr(inErrc, errc))
	for v := range out {
		fmt.Println(v)
	}
}
```

<details>
<summary><strong>Show v1 examples</strong></summary>

### Read all lines from stdin

```go
package main

import (
	"os"

	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	content := λ.Slurp(os.Stdin).UnwrapString()
	_ = content
}
```

### Read a file and pipe it to stdout

```go
package main

import (
	"os"

	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	λ.Open("lorem_ipsum.txt").Slurp().WriteToWriter(os.Stdout)
}
```

### Read a JSON file into a struct

```go
package main

import (
	"fmt"

	λ "github.com/4thel00z/lambda/v1"
)

type MagicSpell struct {
	Name        string  `json:"name"`
	AttackPower float64 `json:"attack_power"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
}

func main() {
	var m MagicSpell
	λ.Open("magic.json").Slurp().JSON(&m).Catch(λ.Die)
	fmt.Printf("%s (%s) atk=%.2f\n%s\n", m.Name, m.Type, m.AttackPower, m.Description)
}
```

### Simple HTTP request

```go
package main

import (
	"os"

	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	λ.Get("https://example.com").Do().Slurp().WriteToWriter(os.Stdout)
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


### Simple AES-CBC encryption with PKCS7 padding

```go
package main

import (
	"crypto/rand"
	"io"
	"strings"

	λ "github.com/4thel00z/lambda/v1"
)

var (
	loremIpsum = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
no sea takimata sanctus est Lorem ipsum dolor sit amet.`
	loremIpsumReader = strings.NewReader(loremIpsum)
)

func getRandomKey() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		panic(err)
	}
	return key
}

func main() {
	key := getRandomKey()
	if λ.Read(loremIpsumReader).EncryptAES(key).DecryptAES(key).UnwrapString() != loremIpsum {
		panic("encryption and decryption doesn't work")
	}

	// test for random payload and key that enc & decryption works fine
	for i := 0; i < 10; i++ {
		key = getRandomKey()
		o := λ.Read(io.LimitReader(rand.Reader, 1024))
		text := o.UnwrapString()
		if o.EncryptAES(key).DecryptAES(key).UnwrapString() != text {
			panic("encryption and decryption doesn't work")
		}
	}
}
```

## Pubkey cryptography

```go
package main

import (
	"log"

	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	rsa := λ.RSA(4096)
	println(rsa.ToPublicKeyPemString())
	println(rsa.ToPrivateKeyPemString())
	original := "Some important message which needs to stay secret lel"
	secretMsg := rsa.EncryptRSA(original).UnwrapString()
	if original == secretMsg {
		log.Fatalln("This encryption don't work boi!")
	}
	if original != rsa.DecryptRSA(secretMsg).UnwrapString() {
		log.Fatalln("This decryption don't work boi!")
	}

	// Use rsa.UnwrapPrivateKey() and rsa.UnwrapPublicKey() if you want to extract the raw key material
	// You can load it via λ.LoadRSA(pub,priv)
}
```

## How to generate a sha256 checksum

```go
package main

import (
	"bytes"

	λ "github.com/4thel00z/lambda/v1"
)

var (
	loremIpsum = `Lorem ipsum dolor sit amet, consetetur sadipscing elitr,
sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum.
Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat,
sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren,
no sea takimata sanctus est Lorem ipsum dolor sit amet.`
	loremIpsumReader = bytes.NewReader([]byte(loremIpsum))
)

func main() {
	expected := "70026299e7c4b3bf5b6256b2069ae0cbc2e0960cad1acb51208a311f3864d5bd"
	if λ.Read(loremIpsumReader).Checksum().UnwrapChecksum() != expected {
		panic("sha256 of loremIpsum is wrong!")
	}
}
```

## How to render markdown and print it to stdout

```go
package main

import (
	"fmt"

	λ "github.com/4thel00z/lambda/v1"
)

func main() {
	out := λ.Markdown().Render(`# Markdown
This is so awesome

## Why is this section so nice
Really dunno`).UnwrapString()

	fmt.Print(out)
}
```

</details>

### Runnable examples

See the runnable examples in [`example/`](example/). For instance:

```bash
go run ./example/json
go run ./example/http
```

v2 examples live in `v2/example/`:

```bash
go run ./v2/example/json
go run ./v2/example/http
go run ./v2/example/conditionals
```

## Development

```bash
go test ./...
go vet ./...
gofmt -w .
```

## License

GPL-3.0 — see [`COPYING`](COPYING).
