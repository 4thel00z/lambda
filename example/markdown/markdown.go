///bin/sh -c true && exec /usr/bin/env go run "$0" "$@"
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
