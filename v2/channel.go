package v2

import (
	"context"
	"errors"
	"io"
)

type chanConfig struct {
	buffer int
}

// ChanOption configures channel helpers (exporters/transforms).
type ChanOption func(*chanConfig)

var errInvalidBuffer = errors.New("lambda/v2: buffer must be >= 0")

// WithBuffer sets the buffer size of output channels created by channel helpers.
// n must be >= 0.
func WithBuffer(n int) ChanOption {
	return func(c *chanConfig) {
		if c == nil {
			return
		}
		c.buffer = n
	}
}

func chanCfg(opts []ChanOption) (chanConfig, error) {
	cfg := chanConfig{buffer: 0}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.buffer < 0 {
		return chanConfig{}, errInvalidBuffer
	}
	return cfg, nil
}

func closedErrStream[T any](err error) (<-chan T, <-chan error) {
	out := make(chan T)
	close(out)
	errc := make(chan error, 1)
	errc <- err
	close(errc)
	return out, errc
}

func closedErr2[T any](err error) (<-chan T, <-chan T, <-chan error) {
	a := make(chan T)
	b := make(chan T)
	close(a)
	close(b)
	errc := make(chan error, 1)
	errc <- err
	close(errc)
	return a, b, errc
}

// Range emits ints from start to end-1 (like Python range(start, end)).
func Range(ctx context.Context, start, end int, opts ...ChanOption) (<-chan int, <-chan error) {
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[int](err)
	}
	out := make(chan int, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		for i := start; i < end; i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case out <- i:
			}
		}
		errc <- nil
	}()

	return out, errc
}

// RangeN emits ints from 0 to n-1.
func RangeN(ctx context.Context, n int, opts ...ChanOption) (<-chan int, <-chan error) {
	return Range(ctx, 0, n, opts...)
}

// FromSlice emits all values from xs, then closes.
func FromSlice[T any](ctx context.Context, xs []T, opts ...ChanOption) (<-chan T, <-chan error) {
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		for i := range xs {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case out <- xs[i]:
			}
		}
		errc <- nil
	}()

	return out, errc
}

// Repeat emits v indefinitely until ctx is canceled.
func Repeat[T any](ctx context.Context, v T, opts ...ChanOption) (<-chan T, <-chan error) {
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		for {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case out <- v:
			}
		}
	}()

	return out, errc
}

// RepeatN emits v n times (or 0 times if n <= 0), then closes.
func RepeatN[T any](ctx context.Context, v T, n int, opts ...ChanOption) (<-chan T, <-chan error) {
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case out <- v:
			}
		}
		errc <- nil
	}()

	return out, errc
}

// GenFn generates values until it returns ok==false or an error.
type GenFn[T any] func(context.Context) (v T, ok bool, err error)

// Generate calls gen until ok==false or an error is returned.
func Generate[T any](ctx context.Context, gen GenFn[T], opts ...ChanOption) (<-chan T, <-chan error) {
	if gen == nil {
		return closedErrStream[T](ErrNilFunc("Generate"))
	}
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		for {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			default:
			}
			v, ok, err := gen(ctx)
			if err != nil {
				errc <- err
				return
			}
			if !ok {
				errc <- nil
				return
			}
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case out <- v:
			}
		}
	}()

	return out, errc
}

// Take forwards up to n values from in to a new channel, then closes the output.
func Take[T any](ctx context.Context, in <-chan T, n int, opts ...ChanOption) (<-chan T, <-chan error) {
	if in == nil {
		return closedErrStream[T](errNilChan)
	}
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		if n <= 0 {
			errc <- nil
			return
		}
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case v, ok := <-in:
				if !ok {
					errc <- nil
					return
				}
				select {
				case <-ctx.Done():
					errc <- ctx.Err()
					return
				case out <- v:
				}
			}
		}
		errc <- nil
	}()

	return out, errc
}

// Drop skips the first n values from in, then forwards the rest to a new channel.
func Drop[T any](ctx context.Context, in <-chan T, n int, opts ...ChanOption) (<-chan T, <-chan error) {
	if in == nil {
		return closedErrStream[T](errNilChan)
	}
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErrStream[T](err)
	}
	out := make(chan T, cfg.buffer)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		ctx = ensureCtx(ctx)
		toDrop := n
		if toDrop < 0 {
			toDrop = 0
		}
		for toDrop > 0 {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case _, ok := <-in:
				if !ok {
					errc <- nil
					return
				}
				toDrop--
			}
		}
		for {
			select {
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			case v, ok := <-in:
				if !ok {
					errc <- nil
					return
				}
				select {
				case <-ctx.Done():
					errc <- ctx.Err()
					return
				case out <- v:
				}
			}
		}
	}()

	return out, errc
}

// Peek consumes one element from in and returns it as first, and returns a new channel out
// that replays that first element and then forwards the rest from in.
//
// If in is closed immediately, first is Err[T](io.EOF), out is closed, and errc yields io.EOF.
func Peek[T any](ctx context.Context, in <-chan T, opts ...ChanOption) (first Option[T], out <-chan T, errc <-chan error) {
	if in == nil {
		ch := make(chan T)
		close(ch)
		ec := make(chan error, 1)
		ec <- errNilChan
		close(ec)
		return Err[T](errNilChan), ch, ec
	}
	cfg, err := chanCfg(opts)
	if err != nil {
		ch := make(chan T)
		close(ch)
		ec := make(chan error, 1)
		ec <- err
		close(ec)
		return Err[T](err), ch, ec
	}
	ctx = ensureCtx(ctx)

	select {
	case <-ctx.Done():
		ch := make(chan T)
		close(ch)
		ec := make(chan error, 1)
		ec <- ctx.Err()
		close(ec)
		return Err[T](ctx.Err()), ch, ec
	case v, ok := <-in:
		if !ok {
			ch := make(chan T)
			close(ch)
			ec := make(chan error, 1)
			ec <- io.EOF
			close(ec)
			return Err[T](io.EOF), ch, ec
		}

		// Ensure we can put the first value into the output channel without blocking.
		buf := cfg.buffer
		if buf < 1 {
			buf = 1
		}
		outCh := make(chan T, buf)
		ec := make(chan error, 1)
		outCh <- v

		go func() {
			defer close(outCh)
			defer close(ec)

			for {
				select {
				case <-ctx.Done():
					ec <- ctx.Err()
					return
				case vv, ok := <-in:
					if !ok {
						ec <- nil
						return
					}
					select {
					case <-ctx.Done():
						ec <- ctx.Err()
						return
					case outCh <- vv:
					}
				}
			}
		}()

		return Ok(v), outCh, ec
	}
}

// Tee splits a stream into two outputs. It blocks if either output blocks (unless buffered).
func Tee[T any](ctx context.Context, in <-chan T, opts ...ChanOption) (out1 <-chan T, out2 <-chan T, errc <-chan error) {
	if in == nil {
		return closedErr2[T](errNilChan)
	}
	cfg, err := chanCfg(opts)
	if err != nil {
		return closedErr2[T](err)
	}

	ctx = ensureCtx(ctx)

	a := make(chan T, cfg.buffer)
	b := make(chan T, cfg.buffer)
	ec := make(chan error, 1)

	go func() {
		defer close(a)
		defer close(b)
		defer close(ec)

		for {
			select {
			case <-ctx.Done():
				ec <- ctx.Err()
				return
			case v, ok := <-in:
				if !ok {
					ec <- nil
					return
				}
				select {
				case <-ctx.Done():
					ec <- ctx.Err()
					return
				case a <- v:
				}
				select {
				case <-ctx.Done():
					ec <- ctx.Err()
					return
				case b <- v:
				}
			}
		}
	}()

	return a, b, ec
}

// Collect drains in into a slice and returns it as Option.
func Collect[T any](ctx context.Context, in <-chan T, opts ...ChanOption) Option[[]T] {
	_ = opts // reserved for future options (e.g. capacity hints)
	if in == nil {
		return Err[[]T](errNilChan)
	}
	ctx = ensureCtx(ctx)
	out := make([]T, 0)
	for {
		select {
		case <-ctx.Done():
			return Err[[]T](ctx.Err())
		case v, ok := <-in:
			if !ok {
				return Ok(out)
			}
			out = append(out, v)
		}
	}
}

// Drain consumes values from in until it is closed or ctx is canceled.
func Drain[T any](ctx context.Context, in <-chan T) error {
	if in == nil {
		return errNilChan
	}
	ctx = ensureCtx(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, ok := <-in:
			if !ok {
				return nil
			}
		}
	}
}
