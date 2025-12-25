package v2

import (
	"errors"
	"testing"
)

func TestConditional_IfElse(t *testing.T) {
	c := If(HasError[int], func(o Option[int]) Option[int] { return Ok(42) }).
		Else(Identity[int])

	if got := c.Do(Ok(1)).Must(); got != 1 {
		t.Fatalf("expected else branch to keep value, got %d", got)
	}
	if got := c.Do(Err[int](errors.New("x"))).Must(); got != 42 {
		t.Fatalf("expected if branch to replace, got %d", got)
	}
}

func TestConditional_ElifAndReuse(t *testing.T) {
	c := If(func(o Option[string]) bool { return o.IsOk() && o.Must() == "a" }, func(Option[string]) Option[string] { return Ok("A") }).
		Elif(func(o Option[string]) bool { return o.IsOk() && o.Must() == "b" }, func(Option[string]) Option[string] { return Ok("B") }).
		Else(func(Option[string]) Option[string] { return Ok("OTHER") })

	if got := c.Do(Ok("a")).Must(); got != "A" {
		t.Fatalf("got %q", got)
	}
	if got := c.Do(Ok("b")).Must(); got != "B" {
		t.Fatalf("got %q", got)
	}
	if got := c.Do(Ok("c")).Must(); got != "OTHER" {
		t.Fatalf("got %q", got)
	}
	// Reuse with different input
	if got := c.Do(Ok("a")).Must(); got != "A" {
		t.Fatalf("reuse got %q", got)
	}
}

func TestConditional_ReturnAndClearError(t *testing.T) {
	c := If(HasError[int], Return(Ok(7))).Else(Identity[int])

	if got := c.Do(Err[int](errors.New("x"))).Must(); got != 7 {
		t.Fatalf("got %d", got)
	}
}

func TestConditional_ClearError(t *testing.T) {
	c := If(HasError[int], ClearError[int]).Else(Identity[int])

	// ClearError keeps value, drops error
	if got := c.Do(Wrap(9, errors.New("x"))).Must(); got != 9 {
		t.Fatalf("got %d", got)
	}
}
