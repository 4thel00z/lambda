package v2

import "gopkg.in/yaml.v3"

// FromYAML unmarshals YAML bytes into a value of type T.
func FromYAML[T any](b Bytes) Option[T] {
	if b.err != nil {
		return Err[T](b.err)
	}
	var v T
	err := yaml.Unmarshal(b.v, &v)
	return Wrap(v, err)
}

// ToYAML marshals the contained value to YAML.
func (o Option[T]) ToYAML() Bytes {
	if o.err != nil {
		return Bytes{Err[[]byte](o.err)}
	}
	b, err := yaml.Marshal(o.v)
	return Bytes{Wrap(b, err)}
}
