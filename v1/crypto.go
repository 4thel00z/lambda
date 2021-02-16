package v1

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

func AddPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	return append(ciphertext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func RemovePKCS7Padding(origData []byte) []byte {
	length := len(origData)
	return origData[:(length - int(origData[length-1]))]
}

func (o Option) Encrypt(key []byte) Option {
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

func (o Option) Decrypt(key []byte) Option {
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
