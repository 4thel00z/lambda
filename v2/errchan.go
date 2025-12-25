package v2

import "errors"

var errNilErrChan = errors.New("lambda/v2: nil error channel")

// JoinErr reads exactly one error from each err channel and joins the non-nil ones.
//
// This is handy because many v2 channel helpers return (<-chan T, <-chan error).
// If all errors are nil, JoinErr returns nil.
func JoinErr(errcs ...<-chan error) error {
	if len(errcs) == 0 {
		return nil
	}
	errs := make([]error, 0, len(errcs))
	for _, ec := range errcs {
		if ec == nil {
			errs = append(errs, errNilErrChan)
			continue
		}
		if err, ok := <-ec; ok && err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}


