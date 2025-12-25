package v2

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var errNilChan = errors.New("lambda/v2: nil input channel")

// ParMapChan maps values read from in in parallel and sends results to the returned channel.
// Result order is not guaranteed.
//
// The returned err channel yields exactly one error (possibly nil) and then closes.
func ParMapChan[T, U any](ctx context.Context, in <-chan T, f MapFn[T, U], opts ...ParOption) (<-chan U, <-chan error) {
	out := make(chan U)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		if in == nil {
			errc <- errNilChan
			return
		}
		if f == nil {
			errc <- ErrNilFunc("ParMapChan")
			return
		}
		cfg, err := parCfg(opts)
		if err != nil {
			errc <- err
			return
		}

		ctx = ensureCtx(ctx)
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(cfg.concurrency)

		for {
			select {
			case <-gctx.Done():
				// Stop scheduling new work; wait for in-flight workers to exit.
				cerr := gctx.Err()
				err = g.Wait()
				if err == nil && cerr != nil {
					err = cerr
				}
				errc <- err
				return

			case v, ok := <-in:
				if !ok {
					errc <- g.Wait()
					return
				}
				vv := v

				g.Go(func() error {
					select {
					case <-gctx.Done():
						return gctx.Err()
					default:
					}

					u := f(vv)

					select {
					case <-gctx.Done():
						return gctx.Err()
					case out <- u:
						return nil
					}
				})
			}
		}
	}()

	return out, errc
}

// ParTryChan maps values read from in in parallel and sends results to the returned channel.
// Result order is not guaranteed.
//
// If any invocation returns an error, the first error is returned (fail-fast) and work is canceled.
// The returned err channel yields exactly one error (possibly nil) and then closes.
func ParTryChan[T, U any](ctx context.Context, in <-chan T, f TryFn[T, U], opts ...ParOption) (<-chan U, <-chan error) {
	out := make(chan U)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		if in == nil {
			errc <- errNilChan
			return
		}
		if f == nil {
			errc <- ErrNilFunc("ParTryChan")
			return
		}
		cfg, err := parCfg(opts)
		if err != nil {
			errc <- err
			return
		}

		ctx = ensureCtx(ctx)
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(cfg.concurrency)

		for {
			select {
			case <-gctx.Done():
				// Stop scheduling new work; wait for in-flight workers to exit.
				cerr := gctx.Err()
				err = g.Wait()
				if err == nil && cerr != nil {
					err = cerr
				}
				errc <- err
				return

			case v, ok := <-in:
				if !ok {
					errc <- g.Wait()
					return
				}
				vv := v

				g.Go(func() error {
					select {
					case <-gctx.Done():
						return gctx.Err()
					default:
					}

					u, err := f(vv)
					if err != nil {
						return err
					}

					select {
					case <-gctx.Done():
						return gctx.Err()
					case out <- u:
						return nil
					}
				})
			}
		}
	}()

	return out, errc
}
