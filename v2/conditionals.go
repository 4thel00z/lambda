package v2

// Predicate decides whether a Conditional branch should run.
type Predicate[T any] func(o Option[T]) bool

// Transformer transforms an Option, usually preserving the value type.
type Transformer[T any] func(o Option[T]) Option[T]

// Conditional represents a reusable if/elif/else chain.
// It does not contain data; you pass the data via Do.
type Conditional[T any] struct {
	predicate        Predicate[T]
	ifBranch         Transformer[T]
	ifElsePredicates []Predicate[T]
	ifElseBranches   []Transformer[T]
	elseBranch       Transformer[T]
}

// Identity returns its input unchanged.
func Identity[T any](o Option[T]) Option[T] { return o }

// HasNoError reports whether o is Ok.
func HasNoError[T any](o Option[T]) bool { return o.err == nil }

// HasError reports whether o is Err.
func HasError[T any](o Option[T]) bool { return o.err != nil }

// ClearError drops the error while keeping the current value.
func ClearError[T any](o Option[T]) Option[T] { return Option[T]{v: o.v} }

// Return creates a transformer that ignores its input and returns option.
func Return[T any](option Option[T]) Transformer[T] {
	return func(Option[T]) Option[T] { return option }
}

// If creates a Conditional chain starting with a predicate and an if-branch.
// If ifBranch is nil, Identity is used.
func If[T any](predicate Predicate[T], ifBranch Transformer[T]) Conditional[T] {
	if ifBranch == nil {
		ifBranch = Identity[T]
	}
	return Conditional[T]{predicate: predicate, ifBranch: ifBranch}
}

// Elif appends an else-if branch.
func (c Conditional[T]) Elif(p Predicate[T], t Transformer[T]) Conditional[T] {
	dstPredicates := append([]Predicate[T]{}, c.ifElsePredicates...)
	dstBranches := append([]Transformer[T]{}, c.ifElseBranches...)
	dstPredicates = append(dstPredicates, p)
	dstBranches = append(dstBranches, t)
	return Conditional[T]{
		predicate:        c.predicate,
		ifBranch:         c.ifBranch,
		ifElsePredicates: dstPredicates,
		ifElseBranches:   dstBranches,
		elseBranch:       c.elseBranch,
	}
}

// Else sets the else branch.
func (c Conditional[T]) Else(t Transformer[T]) Conditional[T] {
	c.elseBranch = t
	return c
}

// Do executes the conditional chain on o.
func (c Conditional[T]) Do(o Option[T]) Option[T] {
	if c.predicate != nil && c.predicate(o) {
		return c.ifBranch(o)
	}
	for i, p := range c.ifElsePredicates {
		if p != nil && p(o) {
			return c.ifElseBranches[i](o)
		}
	}
	if c.elseBranch != nil {
		return c.elseBranch(o)
	}
	return o
}
