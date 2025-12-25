package v2

import (
	"bytes"
	"testing"
)

func TestBytesStringAndWriteTo(t *testing.T) {
	var buf bytes.Buffer

	o := BytesOf([]byte("hi"))
	if o.String().Must() != "hi" {
		t.Fatalf("String mismatch")
	}
	if string(StrOf("hi").Bytes().Must()) != "hi" {
		t.Fatalf("Bytes mismatch")
	}
	n, err := o.WriteTo(&buf)
	if err != nil || n != 2 || buf.String() != "hi" {
		t.Fatalf("WriteTo mismatch: n=%d err=%v buf=%q", n, err, buf.String())
	}
}
