package v2

import (
	"bytes"
	"testing"
)

func TestAESGCMRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0x01}, 32)
	plain := BytesOf([]byte("secret"))
	ct := plain.EncryptAESGCM(key)
	if ct.IsErr() {
		t.Fatalf("encrypt err: %v", ct.Err())
	}
	pt := ct.DecryptAESGCM(key)
	if pt.IsErr() {
		t.Fatalf("decrypt err: %v", pt.Err())
	}
	if string(pt.Must()) != "secret" {
		t.Fatalf("mismatch: %q", string(pt.Must()))
	}
}

func TestRSARoundTripAndPEM(t *testing.T) {
	kp := RSA(2048)
	if kp.IsErr() {
		t.Fatalf("rsa err: %v", kp.Err())
	}
	msg := []byte("hello")
	ct := kp.EncryptOAEP(msg)
	pt := kp.DecryptOAEP(ct.Must()).Must()
	if string(pt) != "hello" {
		t.Fatalf("mismatch: %q", string(pt))
	}

	pub := kp.PublicKeyPEM().Must()
	priv := kp.PrivateKeyPEM().Must()
	loaded := LoadRSA(pub, priv)
	if loaded.IsErr() {
		t.Fatalf("load err: %v", loaded.Err())
	}
	if loaded.FingerprintSHA256().Must() == "" {
		t.Fatalf("expected fingerprint")
	}
}

func TestSHA256Hex(t *testing.T) {
	h := BytesOf([]byte("abc")).SHA256().Hex().Must()
	if h != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("got %s", h)
	}
}


