package stripe

import "errors"

var ErrNotSupported = errors.New("stripe: operation not supported")

func isSessionFailed(err error) bool {
	return err != nil && !errors.Is(err, ErrNotSupported)
}
