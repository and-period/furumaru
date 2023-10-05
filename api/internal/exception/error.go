package exception

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
)

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrFailedPrecondition = errors.New("failed precondition")
	ErrResourceExhausted  = errors.New("resource exhausted")
	ErrNotImplemented     = errors.New("not implemented")
	ErrInternal           = errors.New("internal error")
	ErrCanceled           = errors.New("canceled")
	ErrUnavailable        = errors.New("unavailable")
	ErrDeadlineExceeded   = errors.New("deadline exceeded")
	ErrOutOfRange         = errors.New("out of range")
	ErrUnknown            = errors.New("unknown")
)

func InternalError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return wrapError(ErrCanceled, err)
	case errors.Is(err, context.DeadlineExceeded):
		return wrapError(ErrDeadlineExceeded, err)
	}

	if err := mediaError(err); err != nil {
		return err
	}
	if err := messengerError(err); err != nil {
		return err
	}
	if err := storeError(err); err != nil {
		return err
	}
	if err := userError(err); err != nil {
		return err
	}
	return wrapError(ErrUnknown, err)
}

func wrapError(code, err error) error {
	return fmt.Errorf("%w: %s", code, err.Error())
}

func mediaError(err error) error {
	switch {
	case errors.Is(err, media.ErrInvalidArgument):
		return wrapError(ErrInvalidArgument, err)
	case errors.Is(err, media.ErrNotFound):
		return wrapError(ErrNotFound, err)
	case errors.Is(err, media.ErrAlreadyExists):
		return wrapError(ErrAlreadyExists, err)
	case errors.Is(err, media.ErrFailedPrecondition):
		return wrapError(ErrFailedPrecondition, err)
	case errors.Is(err, media.ErrCanceled):
		return wrapError(ErrCanceled, err)
	case errors.Is(err, media.ErrDeadlineExceeded):
		return wrapError(ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func messengerError(err error) error {
	switch {
	case errors.Is(err, messenger.ErrInvalidArgument):
		return wrapError(ErrInvalidArgument, err)
	case errors.Is(err, messenger.ErrNotFound):
		return wrapError(ErrNotFound, err)
	case errors.Is(err, messenger.ErrAlreadyExists):
		return wrapError(ErrAlreadyExists, err)
	case errors.Is(err, messenger.ErrForbidden):
		return wrapError(ErrForbidden, err)
	case errors.Is(err, messenger.ErrFailedPrecondition):
		return wrapError(ErrFailedPrecondition, err)
	case errors.Is(err, messenger.ErrCanceled):
		return wrapError(ErrCanceled, err)
	case errors.Is(err, messenger.ErrDeadlineExceeded):
		return wrapError(ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func storeError(err error) error {
	switch {
	case errors.Is(err, store.ErrInvalidArgument):
		return wrapError(ErrInvalidArgument, err)
	case errors.Is(err, store.ErrNotFound):
		return wrapError(ErrNotFound, err)
	case errors.Is(err, store.ErrAlreadyExists):
		return wrapError(ErrAlreadyExists, err)
	case errors.Is(err, store.ErrForbidden):
		return wrapError(ErrForbidden, err)
	case errors.Is(err, store.ErrFailedPrecondition):
		return wrapError(ErrFailedPrecondition, err)
	case errors.Is(err, store.ErrUnavailable):
		return wrapError(ErrUnavailable, err)
	case errors.Is(err, store.ErrCanceled):
		return wrapError(ErrCanceled, err)
	case errors.Is(err, store.ErrDeadlineExceeded):
		return wrapError(ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func userError(err error) error {
	switch {
	case errors.Is(err, user.ErrInvalidArgument):
		return wrapError(ErrInvalidArgument, err)
	case errors.Is(err, user.ErrUnauthenticated):
		return wrapError(ErrUnauthenticated, err)
	case errors.Is(err, user.ErrNotFound):
		return wrapError(ErrNotFound, err)
	case errors.Is(err, user.ErrAlreadyExists):
		return wrapError(ErrAlreadyExists, err)
	case errors.Is(err, user.ErrForbidden):
		return wrapError(ErrForbidden, err)
	case errors.Is(err, user.ErrFailedPrecondition):
		return wrapError(ErrFailedPrecondition, err)
	case errors.Is(err, user.ErrCanceled):
		return wrapError(ErrCanceled, err)
	case errors.Is(err, user.ErrDeadlineExceeded):
		return wrapError(ErrDeadlineExceeded, err)
	default:
		return nil
	}
}
