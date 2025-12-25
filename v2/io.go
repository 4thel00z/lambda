package v2

import (
	"errors"
	"io"
	"os"
)

// Open opens a file for reading.
func Open(path string) ReadCloser {
	f, err := os.Open(path)
	var rc io.ReadCloser = f
	return ReadCloser{Wrap(rc, err)}
}

// Create creates/truncates a file for writing.
func Create(path string) WriteCloser {
	f, err := os.Create(path)
	var wc io.WriteCloser = f
	return WriteCloser{Wrap(wc, err)}
}

// Read wraps an io.Reader for pipeline operations.
func Read(r io.Reader) Reader {
	var rr io.Reader = r
	return Reader{Ok(rr)}
}

// ReadAll reads all content from r (no Close).
func ReadAll(r io.Reader) Bytes {
	b, err := io.ReadAll(r)
	return Bytes{Wrap(b, err)}
}

// Slurp reads all content from r and closes it.
func Slurp(r io.ReadCloser) Bytes {
	b, readErr := io.ReadAll(r)
	closeErr := r.Close()
	return Bytes{Wrap(b, errors.Join(readErr, closeErr))}
}
