package v2

import (
	"context"
	"errors"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

func TestParMapChan_CorrectElements(t *testing.T) {
	t.Parallel()

	in := make(chan int)
	out, errc := ParMapChan(context.Background(), in, MapFn[int, int](func(v int) int {
		if v%7 == 0 {
			time.Sleep(1 * time.Millisecond)
		}
		return v * 2
	}), WithConcurrency(5))

	go func() {
		defer close(in)
		for i := 0; i < 100; i++ {
			in <- i
		}
	}()

	var got []int
	for v := range out {
		got = append(got, v)
	}
	err := <-errc
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(got) != 100 {
		t.Fatalf("len(got)=%d, want 100", len(got))
	}
	sort.Ints(got)
	for i := 0; i < 100; i++ {
		want := i * 2
		if got[i] != want {
			t.Fatalf("got[%d]=%d, want %d", i, got[i], want)
		}
	}
}

func TestParTryChan_Error(t *testing.T) {
	t.Parallel()

	in := make(chan int)
	sentinel := errors.New("boom")
	out, errc := ParTryChan(context.Background(), in, TryFn[int, int](func(v int) (int, error) {
		if v == 13 {
			return 0, sentinel
		}
		time.Sleep(1 * time.Millisecond)
		return v, nil
	}), WithConcurrency(8))

	go func() {
		defer close(in)
		for i := 0; i < 50; i++ {
			in <- i
		}
	}()

	// Drain output until it closes.
	for range out {
	}
	err := <-errc
	if !errors.Is(err, sentinel) {
		t.Fatalf("err=%v, want %v", err, sentinel)
	}
}

func TestParMapChan_ConcurrencyLimit(t *testing.T) {
	t.Parallel()

	in := make(chan int)

	var inFlight int64
	var maxInFlight int64

	out, errc := ParMapChan(context.Background(), in, MapFn[int, int](func(v int) int {
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
	}), WithConcurrency(3))

	go func() {
		defer close(in)
		for i := 0; i < 25; i++ {
			in <- i
		}
	}()

	for range out {
	}
	err := <-errc
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if maxInFlight > 3 {
		t.Fatalf("maxInFlight=%d, want <= 3", maxInFlight)
	}
}

func TestParTryChan_ConcurrencyLimit(t *testing.T) {
	t.Parallel()

	in := make(chan int)

	var inFlight int64
	var maxInFlight int64

	out, errc := ParTryChan(context.Background(), in, TryFn[int, int](func(v int) (int, error) {
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
		return v, nil
	}), WithConcurrency(4))

	go func() {
		defer close(in)
		for i := 0; i < 25; i++ {
			in <- i
		}
	}()

	for range out {
	}
	err := <-errc
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if maxInFlight > 4 {
		t.Fatalf("maxInFlight=%d, want <= 4", maxInFlight)
	}
}

func TestParTryChan_CanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := make(chan int)
	out, errc := ParTryChan(ctx, in, TryFn[int, int](func(v int) (int, error) {
		time.Sleep(5 * time.Millisecond)
		return v, nil
	}), WithConcurrency(2))

	// out should close promptly; drain defensively.
	for range out {
	}
	err := <-errc
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context.Canceled", err)
	}
}

func TestParTryChan_NilChan(t *testing.T) {
	t.Parallel()
	out, errc := ParTryChan[int, int](context.Background(), nil, TryFn[int, int](func(v int) (int, error) { return v, nil }))
	for range out {
	}
	err := <-errc
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), "lambda/v2: nil input channel"; got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}
