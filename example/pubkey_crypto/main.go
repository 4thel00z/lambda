package main

import (
	λ "github.com/4thel00z/lambda/v1"
	"log"
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
}
