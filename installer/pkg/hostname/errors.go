package hostname

import "fmt"

type HostnameError struct {
	Err error
}

func (e HostnameError) Error() string {
	return fmt.Sprintf("Hostname error: error=%s", e.Err.Error())
}

func (e HostnameError) Unwrap() error {
	return e.Err
}
