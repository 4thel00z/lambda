package v2

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

// MapFn transforms a T into a U.
type MapFn[T, U any] func(T) U

// TryFn transforms a T into a U, possibly returning an error.
type TryFn[T, U any] func(T) (U, error)

type parConfig struct {
	concurrency int
}

// ParOption configures parallel helpers like ParMap/ParTry.
type ParOption func(*parConfig)

// WithConcurrency sets the maximum number of concurrent workers.
// n must be >= 1.
func WithConcurrency(n int) ParOption {
	return func(c *parConfig) {
		if c == nil {
			return
		}
		c.concurrency = n
	}
}

var errInvalidConcurrency = errors.New("lambda/v2: concurrency must be >= 1")

func parCfg(opts []ParOption) (parConfig, error) {
	cfg := parConfig{concurrency: runtime.GOMAXPROCS(0)}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.concurrency < 1 {
		return parConfig{}, errInvalidConcurrency
	}
	return cfg, nil
}

func ensureCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// ParMap maps each element in parallel. Result order matches input order.
func ParMap[T, U any](ctx context.Context, in []T, f MapFn[T, U], opts ...ParOption) Option[[]U] {
	if f == nil {
		return Err[[]U](ErrNilFunc("ParMap"))
	}
	cfg, err := parCfg(opts)
	if err != nil {
		return Err[[]U](err)
	}
	ctx = ensureCtx(ctx)

	if in == nil {
		return Ok([]U(nil))
	}
	out := make([]U, len(in))

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.concurrency)
	for i := range in {
		i := i
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			v := f(in[i])

			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			out[i] = v
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return Err[[]U](err)
	}
	return Ok(out)
}

// ParTry maps each element in parallel. Result order matches input order.
// If any call returns an error, the first error is returned and the context is canceled.
func ParTry[T, U any](ctx context.Context, in []T, f TryFn[T, U], opts ...ParOption) Option[[]U] {
	if f == nil {
		return Err[[]U](ErrNilFunc("ParTry"))
	}
	cfg, err := parCfg(opts)
	if err != nil {
		return Err[[]U](err)
	}
	ctx = ensureCtx(ctx)

	if in == nil {
		return Ok([]U(nil))
	}
	out := make([]U, len(in))

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.concurrency)
	for i := range in {
		i := i
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			v, err := f(in[i])
			if err != nil {
				return err
			}

			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			out[i] = v
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return Err[[]U](err)
	}
	return Ok(out)
}

// ParMapUnordered maps each element in parallel. Result order is not guaranteed.
func ParMapUnordered[T, U any](ctx context.Context, in []T, f MapFn[T, U], opts ...ParOption) Option[[]U] {
	if f == nil {
		return Err[[]U](ErrNilFunc("ParMapUnordered"))
	}
	cfg, err := parCfg(opts)
	if err != nil {
		return Err[[]U](err)
	}
	ctx = ensureCtx(ctx)

	if in == nil {
		return Ok([]U(nil))
	}

	var (
		mu  sync.Mutex
		out = make([]U, 0, len(in))
	)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.concurrency)
	for i := range in {
		i := i
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			v := f(in[i])

			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			mu.Lock()
			out = append(out, v)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return Err[[]U](err)
	}
	return Ok(out)
}

// ParTryUnordered maps each element in parallel. Result order is not guaranteed.
// If any call returns an error, the first error is returned and the context is canceled.
func ParTryUnordered[T, U any](ctx context.Context, in []T, f TryFn[T, U], opts ...ParOption) Option[[]U] {
	if f == nil {
		return Err[[]U](ErrNilFunc("ParTryUnordered"))
	}
	cfg, err := parCfg(opts)
	if err != nil {
		return Err[[]U](err)
	}
	ctx = ensureCtx(ctx)

	if in == nil {
		return Ok([]U(nil))
	}

	var (
		mu  sync.Mutex
		out = make([]U, 0, len(in))
	)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(cfg.concurrency)
	for i := range in {
		i := i
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			v, err := f(in[i])
			if err != nil {
				return err
			}

			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}

			mu.Lock()
			out = append(out, v)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return Err[[]U](err)
	}
	return Ok(out)
}
