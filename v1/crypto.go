package v1

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"log"
)

func AddPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	return append(ciphertext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func RemovePKCS7Padding(origData []byte) []byte {
	length := len(origData)
	return origData[:(length - int(origData[length-1]))]
}

func (o Option) EncryptAES(key []byte) Option {
	block, err := aes.NewCipher(key)
	if err != nil {
		return Wrap(o.value, err)
	}

	blockSize := block.BlockSize()
	src := AddPKCS7Padding(o.UnwrapBytes(), blockSize)
	dst := make([]byte, blockSize+len(src))
	iv := dst[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return Wrap(o.value, err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(dst[blockSize:], src)

	return Option{
		value: dst,
		err:   nil,
	}
}

func (o Option) DecryptAES(key []byte) Option {
	block, err := aes.NewCipher(key)
	if err != nil {
		return Wrap(o.value, err)
	}
	src := o.UnwrapBytes()
	srcCopy := make([]byte, len(src))
	copy(srcCopy, src)

	blockSize := block.BlockSize()

	if len(srcCopy) < blockSize {
		return Wrap(o.value, errors.New("ciphertext too short"))
	}
	iv := srcCopy[:blockSize]
	srcCopy = srcCopy[blockSize:]
	if len(srcCopy)%blockSize != 0 {
		return Wrap(o.value, errors.New("ciphertext is not a multiple of the block size"))
	}
	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(srcCopy, srcCopy)
	return Option{
		value: RemovePKCS7Padding(srcCopy),
		err:   nil,
	}
}

func (o Option) Checksum() Option {
	return Wrap(sha256.Sum256(o.UnwrapBytes()), nil)
}

func (o Option) Checksum224() Option {
	return Wrap(sha256.Sum224(o.UnwrapBytes()), nil)
}

// LoadRSA parses a public and private pem
func LoadRSA(pub, private []byte) Option {
	pubKey, err := BytesToPublicKey(pub)
	if err != nil {
		return WrapError(err)
	}

	privateKey, err := BytesToPrivateKey(private)
	if err != nil {
		return WrapError(err)
	}

	return WrapValue(RSAKeyPair{
		Private: privateKey,
		Public:  pubKey,
	})
}

type RSAKeyPair struct {
	Private *rsa.PrivateKey `json:"private"`
	Public  *rsa.PublicKey `json:"public"`
}

// RSA generates a new key pair
func RSA(bits int) Option {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return Wrap(RSAKeyPair{}, err)
	}
	return Wrap(RSAKeyPair{Private: privateKey, Public: &privateKey.PublicKey}, nil)
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(key *rsa.PrivateKey) ([]byte, error) {
	if key == nil {
		return nil, errors.New("private key is nil")
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	), nil
}

func (o Option) UnwrapPublicKey() []byte {
	if o.err != nil {
		log.Fatal(o.err)
	}
	keyPair, ok := o.value.(RSAKeyPair)
	if !ok {
		log.Fatal(errors.New("could not UnwrapPublicKeyToBytes, because of type mismatch"))
	}
	memory, err := PublicKeyToBytes(keyPair.Public)
	if err != nil {
		log.Fatal(o.err)
	}
	return memory
}

func (o Option) UnwrapPrivateKey() []byte {
	if o.err != nil {
		log.Fatal(o.err)
	}
	keyPair, ok := o.value.(RSAKeyPair)
	if !ok {
		log.Fatal(errors.New("could not UnwrapPrivateKeyToBytes, because of type mismatch"))
	}
	memory, err := PrivateKeyToBytes(keyPair.Private)
	if err != nil {
		log.Fatal(o.err)
	}
	return memory
}

func (o Option) ToPublicKeyPemString() string {
	return string(o.UnwrapPublicKey())
}

func (o Option) ToPrivateKeyPemString() string {
	return string(o.UnwrapPrivateKey())
}

func (o Option) EncryptRSABytes(data []byte) Option {
	if o.err != nil {
		return o
	}
	return Wrap(EncryptWithPublicKey(data, o.value.(RSAKeyPair).Public))
}

func (o Option) EncryptRSA(data string) Option {
	return o.EncryptRSABytes([]byte(data))
}

func (o Option) DecryptRSABytes(data []byte) Option {
	if o.err != nil {
		return o
	}
	return Wrap(DecryptWithPrivateKey(data, o.value.(RSAKeyPair).Private))
}

func (o Option) DecryptRSA(data string) Option {
	return o.DecryptRSABytes([]byte(data))
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(key *rsa.PublicKey) ([]byte, error) {
	if key == nil {
		return nil, errors.New("public key is nil")
	}
	pubASN1, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	}), nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(memory []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(memory)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err

	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("rsa.PublicKey is not okay")

	}
	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, key *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, key, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
