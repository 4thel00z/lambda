package v1

type Predicate func(o Option) bool

type Conditional struct {
	predicate        Predicate
	ifBranch         Transformer
	ifElsePredicates []Predicate
	ifElseBranches   []Transformer
	elseBranch       Transformer
}

func Identity(option Option) Option {
	return option
}

func HasNoError(o Option) bool {
	return o.err == nil
}

func HasError(o Option) bool {
	return o.err != nil
}

func ClearError(o Option) Option {
	return Option{value: o.value}
}

func If(predicate Predicate, ifBranch Transformer) Conditional {
	if ifBranch == nil {
		ifBranch = Identity
	}
	return Conditional{
		predicate: predicate,
		ifBranch:  ifBranch,
	}
}

func (c Conditional) Elif(p Predicate, i Transformer) Conditional {
	var (
		dstTransformers []Transformer
		dstPredicates   []Predicate
	)
	copy(dstPredicates, c.ifElsePredicates)
	copy(dstTransformers, c.ifElseBranches)
	dstPredicates = append(dstPredicates, p)
	dstTransformers = append(dstTransformers, i)

	return Conditional{
		predicate:        c.predicate,
		ifBranch:         c.ifBranch,
		ifElsePredicates: dstPredicates,
		ifElseBranches:   dstTransformers,
		elseBranch:       c.elseBranch,
	}
}

func (c Conditional) Else(t Transformer) Conditional {
	return Conditional{
		predicate:        c.predicate,
		ifBranch:         c.ifBranch,
		ifElsePredicates: c.ifElsePredicates,
		ifElseBranches:   c.ifElseBranches,
		elseBranch:       t,
	}
}

func (c Conditional) Do(o Option) Option {

	if c.predicate(o) {
		return c.ifBranch(o)
	}
	for i, p := range c.ifElsePredicates {
		if p(o) {
			return c.ifElseBranches[i](o)
		}
	}
	if c.elseBranch != nil {
		return c.elseBranch(o)
	}
	return o
}
