package exception

import "errors"

var (
	ErrInvalidArgument     = errors.New("invalid argument")
	ErrUnauthenticated     = errors.New("unauthenticated")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrFailedPrecondition  = errors.New("failed precondition")
	ErrResourceExhausted   = errors.New("resource exhausted")
	ErrNotImplemented      = errors.New("not implemented")
	ErrInternal            = errors.New("internal error")
	ErrCanceled            = errors.New("canceled")
	ErrUnavailable         = errors.New("unavailable")
	ErrDeadlineExceeded    = errors.New("deadline exceeded")
	ErrOutOfRange          = errors.New("out of range")
	ErrUnknown             = errors.New("unknown")
)

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, ErrResourceExhausted),
		errors.Is(err, ErrUnavailable),
		errors.Is(err, ErrDeadlineExceeded),
		errors.Is(err, ErrInternal):
		return true
	default:
		return false
	}
}
