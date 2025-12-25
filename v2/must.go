package v2

// Must panics if err is non-nil.
//
// This is a small convenience for patterns like:
//
//	λ.Must(λ.JoinErr(errc1, errc2))
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
