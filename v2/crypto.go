package v2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

// SHA256Sum is a pipeline wrapper around Option[[sha256.Size]byte].
//
// This wrapper exists because Go does not allow methods on an instantiated
// generic type whose type argument is a composite type like `[32]byte`.
type SHA256Sum struct{ Option[[sha256.Size]byte] }

// SHA256 returns the SHA256 checksum of the contained bytes.
func (b Bytes) SHA256() SHA256Sum {
	if b.err != nil {
		return SHA256Sum{Err[[sha256.Size]byte](b.err)}
	}
	return SHA256Sum{Ok(sha256.Sum256(b.v))}
}

// Hex formats a sha256 sum as a hex string.
func (s SHA256Sum) Hex() Str {
	if s.err != nil {
		return Str{Err[string](s.err)}
	}
	return Str{Ok(fmt.Sprintf("%x", s.v))}
}

// EncryptAESGCM encrypts bytes using AES-GCM. Output is nonce||ciphertext.
func (b Bytes) EncryptAESGCM(key []byte) Bytes {
	if b.err != nil {
		return Bytes{Err[[]byte](b.err)}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return Bytes{Err[[]byte](err)}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Bytes{Err[[]byte](err)}
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return Bytes{Err[[]byte](err)}
	}
	ct := gcm.Seal(nil, nonce, b.v, nil)
	out := append(nonce, ct...)
	return Bytes{Ok(out)}
}

// DecryptAESGCM decrypts bytes produced by EncryptAESGCM (nonce||ciphertext).
func (b Bytes) DecryptAESGCM(key []byte) Bytes {
	if b.err != nil {
		return Bytes{Err[[]byte](b.err)}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return Bytes{Err[[]byte](err)}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Bytes{Err[[]byte](err)}
	}
	ns := gcm.NonceSize()
	if len(b.v) < ns {
		return Bytes{Err[[]byte](errors.New("lambda/v2: ciphertext too short"))}
	}
	nonce := b.v[:ns]
	ct := b.v[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	return Bytes{Wrap(pt, err)}
}

// RSAKeyPair holds an RSA private/public key pair.
type RSAKeyPair struct {
	Private *rsa.PrivateKey
	Public  *rsa.PublicKey
}

// RSAKeys is a pipeline wrapper around Option[RSAKeyPair].
type RSAKeys struct{ Option[RSAKeyPair] }

// RSA generates a new RSA key pair.
func RSA(bits int) RSAKeys {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return RSAKeys{Err[RSAKeyPair](err)}
	}
	return RSAKeys{Ok(RSAKeyPair{Private: priv, Public: &priv.PublicKey})}
}

// EncryptOAEP encrypts plaintext using the keypair's public key (SHA256 OAEP).
func (k RSAKeys) EncryptOAEP(plaintext []byte) Bytes {
	if k.err != nil {
		return Bytes{Err[[]byte](k.err)}
	}
	if k.v.Public == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil public key"))}
	}
	ct, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, k.v.Public, plaintext, nil)
	return Bytes{Wrap(ct, err)}
}

// DecryptOAEP decrypts ciphertext using the keypair's private key (SHA256 OAEP).
func (k RSAKeys) DecryptOAEP(ciphertext []byte) Bytes {
	if k.err != nil {
		return Bytes{Err[[]byte](k.err)}
	}
	if k.v.Private == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil private key"))}
	}
	pt, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, k.v.Private, ciphertext, nil)
	return Bytes{Wrap(pt, err)}
}

// PublicKeyPEM encodes the public key to PEM (PKIX).
func (k RSAKeys) PublicKeyPEM() Bytes {
	if k.err != nil {
		return Bytes{Err[[]byte](k.err)}
	}
	if k.v.Public == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil public key"))}
	}
	asn1, err := x509.MarshalPKIXPublicKey(k.v.Public)
	if err != nil {
		return Bytes{Err[[]byte](err)}
	}
	return Bytes{Ok(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: asn1}))}
}

// PrivateKeyPEM encodes the private key to PEM (PKCS#1).
func (k RSAKeys) PrivateKeyPEM() Bytes {
	if k.err != nil {
		return Bytes{Err[[]byte](k.err)}
	}
	if k.v.Private == nil {
		return Bytes{Err[[]byte](errors.New("lambda/v2: nil private key"))}
	}
	return Bytes{Ok(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k.v.Private)}))}
}

// LoadRSA loads an RSA key pair from PEM-encoded public/private keys.
func LoadRSA(pubPEM, privPEM []byte) RSAKeys {
	pub, err := ParsePublicKeyPEM(pubPEM)
	if err != nil {
		return RSAKeys{Err[RSAKeyPair](err)}
	}
	priv, err := ParsePrivateKeyPEM(privPEM)
	if err != nil {
		return RSAKeys{Err[RSAKeyPair](err)}
	}
	return RSAKeys{Ok(RSAKeyPair{Private: priv, Public: pub})}
}

// ParsePublicKeyPEM parses a PKIX public key PEM into *rsa.PublicKey.
func ParsePublicKeyPEM(b []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("lambda/v2: invalid public key pem")
	}
	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("lambda/v2: not an rsa public key")
	}
	return pub, nil
}

// ParsePrivateKeyPEM parses an RSA private key PEM (PKCS#1) into *rsa.PrivateKey.
func ParsePrivateKeyPEM(b []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("lambda/v2: invalid private key pem")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// FingerprintSHA256 computes SSH-style base64 SHA256 fingerprint of the public key.
func (k RSAKeys) FingerprintSHA256() Str {
	if k.err != nil {
		return Str{Err[string](k.err)}
	}
	pubPem := k.PublicKeyPEM()
	if pubPem.err != nil {
		return Str{Err[string](pubPem.err)}
	}
	sum := sha256.Sum256(pubPem.v)
	return Str{Ok("SHA256:" + base64.RawStdEncoding.EncodeToString(sum[:]))}
}
