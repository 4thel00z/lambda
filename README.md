# λ

![lambda-tests](https://github.com/4thel00z/lambda/workflows/Test/badge.svg)
![lambda-logo](https://raw.githubusercontent.com/4thel00z/lambda/master/logo.png)

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
	"fmt"
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

    fmt.Println(strings.Join([]string{m.Name, m.Type, fmt.Sprintf("%f", m.AttackPower), m.Description}, "\n"))

	// ToJSON() detects if the current value is a pointer or not
	fmt.Println(λ.WrapValue(m).ToJSON().UnwrapString())
	// Works even if you use the pointer operator again
	fmt.Println(λ.WrapValue(&m).ToJSON().UnwrapString())
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

### Simple AES-CBC encryption with PKC7 Padding

```go
package main

import (
	"crypto/rand"
	λ "github.com/4thel00z/lambda/v1"
	"io"
	"strings"
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
	if λ.Read(loremIpsumReader).Encrypt(key).Decrypt(key).UnwrapString() != loremIpsum {
		panic("encryption and decryption doesn't work")
	}

	// test for random payload and key that enc & decryption works fine
	for i := 0; i < 10; i++ {
		key = getRandomKey()
		o := λ.Read(io.LimitReader(rand.Reader, 1024))
		text := o.UnwrapString()
		if o.Encrypt(key).Decrypt(key).UnwrapString() != text {
			panic("encryption and decryption doesn't work")
		}
	}

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

## How to render a markdown and print out to stdout

```go
package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"os"
)

func main() {
	λ.Markdown().Render(`# Markdown
This is so awesome

## Why is this section so nice
Really dunno

### Omg, can do all the things
* yeah
* all
* of
* them

#### Emojis work too:👩
#### Code aswell:
`+"```"+ `
import "github.com/charmbracelet/glamour"

r, _ := glamour.NewTermRenderer(
    // detect background color and pick either the default dark or light theme
    glamour.WithAutoStyle(),
    // wrap output at specific width
    glamour.WithWordWrap(40),
)

out, err := r.Render(in)
fmt.Print(out)
`+"```").WriteStringTo(os.Stdin)
}
```

## Todo

* Make Option more flexible an pretty

## License

This project is licensed under the GPL-3 license.
