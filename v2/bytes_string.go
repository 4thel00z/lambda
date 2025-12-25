package v2

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// Bytes is a pipeline wrapper around Option[[]byte].
type Bytes struct{ Option[[]byte] }

// Str is a pipeline wrapper around Option[string].
type Str struct{ Option[string] }

// Lines is a pipeline wrapper around Option[[]string].
type Lines struct{ Option[[]string] }

// Reader is a pipeline wrapper around Option[io.Reader].
type Reader struct{ Option[io.Reader] }

// ReadCloser is a pipeline wrapper around Option[io.ReadCloser].
type ReadCloser struct{ Option[io.ReadCloser] }

// WriteCloser is a pipeline wrapper around Option[io.WriteCloser].
type WriteCloser struct{ Option[io.WriteCloser] }

// BytesOf wraps a raw []byte into a Bytes pipeline.
func BytesOf(v []byte) Bytes { return Bytes{Ok(v)} }

// StrOf wraps a raw string into a Str pipeline.
func StrOf(v string) Str { return Str{Ok(v)} }

// LinesOf wraps a raw []string into a Lines pipeline.
func LinesOf(v []string) Lines { return Lines{Ok(v)} }

// String converts bytes to string.
func (b Bytes) String() Str {
	if b.err != nil {
		return Str{Err[string](b.err)}
	}
	return Str{Ok(string(b.v))}
}

// Bytes converts string to bytes.
func (s Str) Bytes() Bytes {
	if s.err != nil {
		return Bytes{Err[[]byte](s.err)}
	}
	return Bytes{Ok([]byte(s.v))}
}

// Reader converts bytes to an io.Reader.
func (b Bytes) Reader() Reader {
	if b.err != nil {
		return Reader{Err[io.Reader](b.err)}
	}
	var r io.Reader = bytes.NewReader(b.v)
	return Reader{Ok(r)}
}

// Reader converts string to an io.Reader.
func (s Str) Reader() Reader {
	if s.err != nil {
		return Reader{Err[io.Reader](s.err)}
	}
	var r io.Reader = strings.NewReader(s.v)
	return Reader{Ok(r)}
}

// ReadAll reads all content from the contained io.Reader.
func (r Reader) ReadAll() Bytes {
	if r.err != nil {
		return Bytes{Err[[]byte](r.err)}
	}
	return ReadAll(r.v)
}

// Slurp reads all content from the contained io.ReadCloser and closes it.
func (r ReadCloser) Slurp() Bytes {
	if r.err != nil {
		return Bytes{Err[[]byte](r.err)}
	}
	return Slurp(r.v)
}

// WriteTo implements io.WriterTo for Bytes.
func (b Bytes) WriteTo(w io.Writer) (int64, error) {
	if b.err != nil {
		return 0, b.err
	}
	if w == nil {
		return 0, errors.New("lambda/v2: nil writer")
	}
	n, err := w.Write(b.v)
	return int64(n), err
}

// WriteToWriter preserves chain-friendly behavior.
func (b Bytes) WriteToWriter(w io.Writer) Bytes {
	_, err := b.WriteTo(w)
	return Bytes{Wrap(b.v, err)}
}
