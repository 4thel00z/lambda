package v2

import (
	"errors"
	"testing"
)

func TestOption_OkGetMust(t *testing.T) {
	o := Ok(123)
	if o.IsErr() || !o.IsOk() {
		t.Fatalf("expected ok")
	}
	v, err := o.Get()
	if err != nil || v != 123 {
		t.Fatalf("Get() got (%v, %v)", v, err)
	}
	if o.Must() != 123 {
		t.Fatalf("Must() mismatch")
	}
}

func TestOption_ErrGetMustPanics(t *testing.T) {
	sentinel := errors.New("boom")
	o := Err[int](sentinel)
	if !o.IsErr() || o.IsOk() {
		t.Fatalf("expected err")
	}
	_, err := o.Get()
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	_ = o.Must()
}

func TestOption_CatchTapTapErr(t *testing.T) {
	var tapped int
	var tappedErr error

	Ok(5).
		Tap(func(v int) { tapped = v }).
		TapErr(func(err error) { tappedErr = err })

	if tapped != 5 {
		t.Fatalf("Tap not called")
	}
	if tappedErr != nil {
		t.Fatalf("TapErr should not be called")
	}

	sentinel := errors.New("x")
	tapped = 0
	tappedErr = nil

	Err[int](sentinel).
		Tap(func(v int) { tapped = v }).
		TapErr(func(err error) { tappedErr = err }).
		Catch(func(err error) error { return errors.Join(err, errors.New("y")) })

	if tapped != 0 {
		t.Fatalf("Tap should not be called on Err")
	}
	if tappedErr == nil {
		t.Fatalf("TapErr should be called on Err")
	}
}

func TestMapThenTry(t *testing.T) {
	o := Ok(2)
	m := Map(o, func(v int) string { return "n=" + string(rune('0'+v)) })
	if m.IsErr() {
		t.Fatalf("unexpected err: %v", m.Err())
	}
	if m.Must() == "" {
		t.Fatalf("expected non-empty")
	}

	t1 := Then(Ok("x"), func(s string) Option[int] { return Ok(len(s)) })
	if t1.Must() != 1 {
		t.Fatalf("Then result mismatch")
	}

	sentinel := errors.New("nope")
	t2 := Try(Ok("x"), func(string) (int, error) { return 0, sentinel })
	if !errors.Is(t2.Err(), sentinel) {
		t.Fatalf("expected propagated error")
	}
}


