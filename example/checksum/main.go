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
