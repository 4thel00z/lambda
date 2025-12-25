package v2

import (
	"context"
	"errors"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

func TestParMap_Ordered(t *testing.T) {
	t.Parallel()

	in := make([]int, 100)
	for i := range in {
		in[i] = i
	}

	out := ParMap(context.Background(), in, MapFn[int, int](func(v int) int {
		// Stagger to increase the chance of reordering if indexing is wrong.
		if v%7 == 0 {
			time.Sleep(2 * time.Millisecond)
		}
		return v * v
	}), WithConcurrency(8)).Must()

	if len(out) != len(in) {
		t.Fatalf("len(out)=%d, want %d", len(out), len(in))
	}
	for i := range in {
		want := in[i] * in[i]
		if out[i] != want {
			t.Fatalf("out[%d]=%d, want %d", i, out[i], want)
		}
	}
}

func TestParTry_Ordered(t *testing.T) {
	t.Parallel()

	in := []int{0, 1, 2, 3, 4, 5}
	out, err := ParTry(context.Background(), in, TryFn[int, int](func(v int) (int, error) {
		if v%2 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
		return v + 10, nil
	}), WithConcurrency(3)).Get()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []int{10, 11, 12, 13, 14, 15}
	for i := range want {
		if out[i] != want[i] {
			t.Fatalf("out[%d]=%d, want %d", i, out[i], want[i])
		}
	}
}

func TestParTry_Errors(t *testing.T) {
	t.Parallel()

	in := []int{0, 1, 2, 3, 4}
	sentinel := errors.New("boom")
	_, err := ParTry(context.Background(), in, TryFn[int, int](func(v int) (int, error) {
		if v == 3 {
			return 0, sentinel
		}
		return v, nil
	}), WithConcurrency(4)).Get()
	if !errors.Is(err, sentinel) {
		t.Fatalf("err=%v, want %v", err, sentinel)
	}
}

func TestParMapUnordered_CorrectElements(t *testing.T) {
	t.Parallel()

	in := make([]int, 200)
	for i := range in {
		in[i] = i
	}

	out, err := ParMapUnordered(context.Background(), in, MapFn[int, int](func(v int) int {
		// Add jitter to scramble completion order.
		if v%5 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
		return v
	}), WithConcurrency(12)).Get()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(out) != len(in) {
		t.Fatalf("len(out)=%d, want %d", len(out), len(in))
	}

	sort.Ints(out)
	for i := range in {
		if out[i] != in[i] {
			t.Fatalf("out[%d]=%d, want %d", i, out[i], in[i])
		}
	}
}

func TestParMap_ConcurrencyLimit(t *testing.T) {
	t.Parallel()

	in := make([]int, 50)
	for i := range in {
		in[i] = i
	}

	var inFlight int64
	var maxInFlight int64

	_, err := ParMap(context.Background(), in, MapFn[int, int](func(v int) int {
		n := atomic.AddInt64(&inFlight, 1)
		for {
			m := atomic.LoadInt64(&maxInFlight)
			if n <= m {
				break
			}
			if atomic.CompareAndSwapInt64(&maxInFlight, m, n) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		atomic.AddInt64(&inFlight, -1)
		return v
	}), WithConcurrency(3)).Get()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if maxInFlight > 3 {
		t.Fatalf("maxInFlight=%d, want <= 3", maxInFlight)
	}
}

func TestParMap_NilFunc(t *testing.T) {
	t.Parallel()
	_, err := ParMap(context.Background(), []int{1}, MapFn[int, int](nil)).Get()
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), (ErrNilFunc("ParMap")).Error(); got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}

func TestParMap_InvalidConcurrency(t *testing.T) {
	t.Parallel()
	_, err := ParMap(context.Background(), []int{1}, MapFn[int, int](func(v int) int { return v }), WithConcurrency(0)).Get()
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), "lambda/v2: concurrency must be >= 1"; got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}

func TestParMap_CanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := ParMap(ctx, []int{1, 2, 3}, MapFn[int, int](func(v int) int {
		time.Sleep(5 * time.Millisecond)
		return v
	}), WithConcurrency(2)).Get()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context.Canceled", err)
	}
}
