package locale

import "fmt"

// LocaleGenError represents an error that occured
// when setting up or validating the new install locales.
type LocaleGenError struct {
	Err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e LocaleGenError) Error() string {
	return fmt.Sprintf("Error while configuring locales: error=%s", e.Err.Error())
}

// Unwrap returns the original error wrapped inside
// LocaleGenError.
func (e LocaleGenError) Unwrap() error {
	return e.Err
}
