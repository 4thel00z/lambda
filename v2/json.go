package v2

import "encoding/json"

// FromJSON unmarshals JSON bytes into a value of type T.
func FromJSON[T any](b Bytes) Option[T] {
	if b.err != nil {
		return Err[T](b.err)
	}
	var v T
	err := json.Unmarshal(b.v, &v)
	return Wrap(v, err)
}

// ToJSON marshals the contained value to JSON.
func (o Option[T]) ToJSON() Bytes {
	if o.err != nil {
		return Bytes{Err[[]byte](o.err)}
	}
	b, err := json.Marshal(o.v)
	return Bytes{Wrap(b, err)}
}
