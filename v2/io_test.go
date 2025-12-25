package v2

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type testReadCloser struct {
	r        *bytes.Reader
	closed   bool
	closeErr error
}

func (t *testReadCloser) Read(p []byte) (int, error) { return t.r.Read(p) }
func (t *testReadCloser) Close() error {
	t.closed = true
	return t.closeErr
}

func TestSlurp_Closes(t *testing.T) {
	rc := &testReadCloser{r: bytes.NewReader([]byte("hello"))}
	o := Slurp(rc)
	if !rc.closed {
		t.Fatalf("expected Close to be called")
	}
	if o.Must() == nil || string(o.Must()) != "hello" {
		t.Fatalf("unexpected content")
	}
}

func TestSlurp_JoinsCloseError(t *testing.T) {
	sentinel := errors.New("close failed")
	rc := &testReadCloser{r: bytes.NewReader([]byte("x")), closeErr: sentinel}
	o := Slurp(rc)
	if !errors.Is(o.Err(), sentinel) {
		t.Fatalf("expected joined close error")
	}
}

func TestOptionReader_ReadAll(t *testing.T) {
	var r io.Reader = bytes.NewReader([]byte("abc"))
	o := Read(r).ReadAll()
	if string(o.Must()) != "abc" {
		t.Fatalf("unexpected")
	}
}
