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
