package v2

import "fmt"

// ErrNilFunc is returned when a required callback is nil.
type ErrNilFunc string

func (e ErrNilFunc) Error() string {
	return fmt.Sprintf("lambda/v2: nil func passed to %s", string(e))
}
