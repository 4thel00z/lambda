package v2

// Option is a value+error container for fluent pipelines.
//
// Important Go generics note:
// Go (as of 1.23) does not allow methods with their own type parameters.
// That means type-changing operations like Map/Then/Try are implemented as
// generic functions (see Map/Then/Try below), while type-preserving
// operations are methods.
type Option[T any] struct {
	v   T
	err error
}

// Ok constructs an Option holding a value with no error.
func Ok[T any](v T) Option[T] { return Option[T]{v: v} }

// Err constructs an Option holding an error (value will be T's zero value).
func Err[T any](err error) Option[T] {
	var z T
	return Option[T]{v: z, err: err}
}

// Wrap constructs an Option from a value and an error.
func Wrap[T any](v T, err error) Option[T] { return Option[T]{v: v, err: err} }

// Get returns the underlying value and error.
func (o Option[T]) Get() (T, error) { return o.v, o.err }

// Must returns the underlying value or panics if the Option contains an error.
func (o Option[T]) Must() T {
	if o.err != nil {
		panic(o.err)
	}
	return o.v
}

// Err returns the underlying error.
func (o Option[T]) Err() error { return o.err }

// IsOk reports whether the Option contains no error.
func (o Option[T]) IsOk() bool { return o.err == nil }

// IsErr reports whether the Option contains an error.
func (o Option[T]) IsErr() bool { return o.err != nil }

// Catch transforms the current error (if any) and returns the updated Option.
func (o Option[T]) Catch(handler func(error) error) Option[T] {
	if o.err == nil || handler == nil {
		return o
	}
	o.err = handler(o.err)
	return o
}

// Tap executes f for side effects if the Option is Ok.
func (o Option[T]) Tap(f func(T)) Option[T] {
	if o.err == nil && f != nil {
		f(o.v)
	}
	return o
}

// TapErr executes f for side effects if the Option is Err.
func (o Option[T]) TapErr(f func(error)) Option[T] {
	if o.err != nil && f != nil {
		f(o.err)
	}
	return o
}

// Or returns the contained value if Ok, otherwise returns fallback.
func (o Option[T]) Or(fallback T) T {
	if o.err != nil {
		return fallback
	}
	return o.v
}

// Map transforms the value using f if o is Ok, otherwise propagates the error.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.err != nil {
		return Err[U](o.err)
	}
	if f == nil {
		return Err[U](ErrNilFunc("Map"))
	}
	return Ok(f(o.v))
}

// Then transforms the value using f if o is Ok, otherwise propagates the error.
// This is the monadic bind / flat-map.
func Then[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.err != nil {
		return Err[U](o.err)
	}
	if f == nil {
		return Err[U](ErrNilFunc("Then"))
	}
	return f(o.v)
}

// Try transforms the value using f if o is Ok, otherwise propagates the error.
func Try[T, U any](o Option[T], f func(T) (U, error)) Option[U] {
	if o.err != nil {
		return Err[U](o.err)
	}
	if f == nil {
		return Err[U](ErrNilFunc("Try"))
	}
	v, err := f(o.v)
	return Wrap(v, err)
}
