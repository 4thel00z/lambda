package v2

import (
	"context"
	"errors"
	"io"
	"sort"
	"testing"
)

func mustErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestRangeN(t *testing.T) {
	t.Parallel()

	ch, errc := RangeN(context.Background(), 5)
	got := Collect[int](context.Background(), ch).Must()
	mustErr(t, <-errc)
	if want := []int{0, 1, 2, 3, 4}; len(got) != len(want) {
		t.Fatalf("len=%d, want %d", len(got), len(want))
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got[%d]=%d, want %d", i, got[i], want[i])
			}
		}
	}
}

func TestFromSlice(t *testing.T) {
	t.Parallel()

	ch, errc := FromSlice(context.Background(), []string{"a", "b", "c"})
	got := Collect[string](context.Background(), ch).Must()
	mustErr(t, <-errc)
	if want := []string{"a", "b", "c"}; len(got) != len(want) {
		t.Fatalf("len=%d, want %d", len(got), len(want))
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got[%d]=%q, want %q", i, got[i], want[i])
			}
		}
	}
}

func TestRepeatN(t *testing.T) {
	t.Parallel()

	ch, errc := RepeatN(context.Background(), "x", 3)
	got := Collect[string](context.Background(), ch).Must()
	mustErr(t, <-errc)
	if want := []string{"x", "x", "x"}; len(got) != len(want) {
		t.Fatalf("len=%d, want %d", len(got), len(want))
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got[%d]=%q, want %q", i, got[i], want[i])
			}
		}
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	i := 0
	ch, errc := Generate(context.Background(), GenFn[int](func(ctx context.Context) (int, bool, error) {
		if i >= 4 {
			return 0, false, nil
		}
		i++
		return i, true, nil
	}))

	got := Collect[int](context.Background(), ch).Must()
	mustErr(t, <-errc)
	if want := []int{1, 2, 3, 4}; len(got) != len(want) {
		t.Fatalf("len=%d, want %d", len(got), len(want))
	} else {
		for idx := range want {
			if got[idx] != want[idx] {
				t.Fatalf("got[%d]=%d, want %d", idx, got[idx], want[idx])
			}
		}
	}
}

func TestTakeAndDrop(t *testing.T) {
	t.Parallel()

	src, srcErrc := RangeN(context.Background(), 10, WithBuffer(10))
	defer func() { _ = <-srcErrc }()

	take, takeErrc := Take(context.Background(), src, 3)
	gotTake := Collect[int](context.Background(), take).Must()
	mustErr(t, <-takeErrc)
	if want := []int{0, 1, 2}; len(gotTake) != len(want) {
		t.Fatalf("take len=%d, want %d", len(gotTake), len(want))
	} else {
		for i := range want {
			if gotTake[i] != want[i] {
				t.Fatalf("take got[%d]=%d, want %d", i, gotTake[i], want[i])
			}
		}
	}

	// New source for drop.
	src2, src2Errc := RangeN(context.Background(), 6, WithBuffer(6))
	defer func() { _ = <-src2Errc }()

	drop, dropErrc := Drop(context.Background(), src2, 2)
	gotDrop := Collect[int](context.Background(), drop).Must()
	mustErr(t, <-dropErrc)
	if want := []int{2, 3, 4, 5}; len(gotDrop) != len(want) {
		t.Fatalf("drop len=%d, want %d", len(gotDrop), len(want))
	} else {
		for i := range want {
			if gotDrop[i] != want[i] {
				t.Fatalf("drop got[%d]=%d, want %d", i, gotDrop[i], want[i])
			}
		}
	}
}

func TestPeek_ReplaysFirst(t *testing.T) {
	t.Parallel()

	src, srcErrc := FromSlice(context.Background(), []int{7, 8, 9})
	defer func() { _ = <-srcErrc }()

	first, out, errc := Peek(context.Background(), src)
	if first.Err() != nil {
		t.Fatalf("first.Err=%v", first.Err())
	}
	if first.Must() != 7 {
		t.Fatalf("first=%d, want 7", first.Must())
	}

	got := Collect[int](context.Background(), out).Must()
	mustErr(t, <-errc)
	if want := []int{7, 8, 9}; len(got) != len(want) {
		t.Fatalf("len=%d, want %d", len(got), len(want))
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("got[%d]=%d, want %d", i, got[i], want[i])
			}
		}
	}
}

func TestPeek_EmptyChannel(t *testing.T) {
	t.Parallel()

	ch := make(chan int)
	close(ch)

	first, out, errc := Peek(context.Background(), ch)
	if !errors.Is(first.Err(), io.EOF) {
		t.Fatalf("first.Err=%v, want io.EOF", first.Err())
	}
	_ = Collect[int](context.Background(), out).Must()
	if err := <-errc; !errors.Is(err, io.EOF) {
		t.Fatalf("err=%v, want io.EOF", err)
	}
}

func TestTee(t *testing.T) {
	t.Parallel()

	src, srcErrc := FromSlice(context.Background(), []int{1, 2, 3, 4})
	defer func() { _ = <-srcErrc }()

	a, b, errc := Tee(context.Background(), src, WithBuffer(4))

	gotA := Collect[int](context.Background(), a).Must()
	gotB := Collect[int](context.Background(), b).Must()
	mustErr(t, <-errc)

	sort.Ints(gotA)
	sort.Ints(gotB)
	if len(gotA) != 4 || len(gotB) != 4 {
		t.Fatalf("lens=%d/%d, want 4/4", len(gotA), len(gotB))
	}
	for i := 0; i < 4; i++ {
		if gotA[i] != i+1 || gotB[i] != i+1 {
			t.Fatalf("tee mismatch at %d: %v %v", i, gotA, gotB)
		}
	}
}

func TestCollect_NilChan(t *testing.T) {
	t.Parallel()
	_, err := Collect[int](context.Background(), nil).Get()
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), "lambda/v2: nil input channel"; got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}

func TestWithBuffer_Invalid(t *testing.T) {
	t.Parallel()
	ch, errc := RangeN(context.Background(), 1, WithBuffer(-1))
	_ = Collect[int](context.Background(), ch).Must() // should be empty
	err := <-errc
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), "lambda/v2: buffer must be >= 0"; got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}

func TestCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch, errc := Repeat(ctx, 1)
	_ = Collect[int](context.Background(), ch) // may be empty
	if err := <-errc; !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context.Canceled", err)
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	in := make(chan int) // no sender; Take should exit on ctx cancel without leaking senders
	out, outErrc := Take(ctx2, in, 10)
	_ = Collect[int](context.Background(), out).Must()
	if err := <-outErrc; !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context.Canceled", err)
	}
}

func TestGenerate_NilFunc(t *testing.T) {
	t.Parallel()
	ch, errc := Generate[int](context.Background(), nil)
	_ = Collect[int](context.Background(), ch).Must()
	err := <-errc
	if err == nil {
		t.Fatalf("expected error")
	}
	if got, want := err.Error(), (ErrNilFunc("Generate")).Error(); got != want {
		t.Fatalf("err=%q, want %q", got, want)
	}
}

func TestDrain(t *testing.T) {
	t.Parallel()
	ch, errc := FromSlice(context.Background(), []int{1, 2, 3}, WithBuffer(3))
	mustErr(t, <-errc)
	if err := Drain(context.Background(), ch); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}
