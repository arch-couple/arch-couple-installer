package mirrors

import "fmt"

// MirrorListError represents an error that occured
// during setting or validating mirrors
type MirrorListError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e MirrorListError) Error() string {
	return fmt.Sprintf("mirrorlist error: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside
// MirrorListError.
func (e MirrorListError) Unwrap() error {
	return e.err
}
