package exception

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/firebase/messaging"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/postalcode"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"gorm.io/gorm"
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

	if isInternal(err) {
		return err
	}

	switch {
	case errors.Is(err, context.Canceled):
		return wrapError("canceled", ErrCanceled, err)
	case errors.Is(err, context.DeadlineExceeded):
		return wrapError("deadline exceeded", ErrDeadlineExceeded, err)
	}

	if err := mediaError(err); err != nil {
		return err
	}
	if err := messengerError(err); err != nil {
		return err
	}
	if err := validationError(err); err != nil {
		return err
	}
	if err := dbError(err); err != nil {
		return err
	}
	if err := authError(err); err != nil {
		return err
	}
	if err := storageError(err); err != nil {
		return err
	}
	if err := mailerError(err); err != nil {
		return err
	}
	if err := messagingError(err); err != nil {
		return err
	}
	if err := notifierError(err); err != nil {
		return err
	}
	if err := externalError(err); err != nil {
		return err
	}
	return wrapError("internal", ErrUnknown, err)
}

func Retryable(err error) bool {
	return errors.Is(err, ErrResourceExhausted) ||
		errors.Is(err, ErrCanceled) ||
		errors.Is(err, ErrUnavailable) ||
		errors.Is(err, ErrDeadlineExceeded) ||
		errors.Is(err, ErrOutOfRange)
}

func wrapError(prefix string, code, err error) error {
	if prefix == "" {
		return fmt.Errorf("%w: %s", code, err.Error())
	}
	return fmt.Errorf("%s: %w: %s", prefix, code, err.Error())
}

func mediaError(err error) error {
	switch {
	case errors.Is(err, media.ErrInvalidArgument):
		return wrapError("", ErrInvalidArgument, err)
	case errors.Is(err, media.ErrNotFound):
		return wrapError("", ErrNotFound, err)
	case errors.Is(err, media.ErrAlreadyExists):
		return wrapError("", ErrAlreadyExists, err)
	case errors.Is(err, media.ErrFailedPrecondition):
		return wrapError("", ErrFailedPrecondition, err)
	case errors.Is(err, media.ErrCanceled):
		return wrapError("", ErrCanceled, err)
	case errors.Is(err, media.ErrDeadlineExceeded):
		return wrapError("", ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func messengerError(err error) error {
	switch {
	case errors.Is(err, messenger.ErrInvalidArgument):
		return wrapError("", ErrInvalidArgument, err)
	case errors.Is(err, messenger.ErrNotFound):
		return wrapError("", ErrNotFound, err)
	case errors.Is(err, messenger.ErrAlreadyExists):
		return wrapError("", ErrAlreadyExists, err)
	case errors.Is(err, messenger.ErrForbidden):
		return wrapError("", ErrForbidden, err)
	case errors.Is(err, messenger.ErrFailedPrecondition):
		return wrapError("", ErrFailedPrecondition, err)
	case errors.Is(err, messenger.ErrCanceled):
		return wrapError("", ErrCanceled, err)
	case errors.Is(err, messenger.ErrDeadlineExceeded):
		return wrapError("", ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func isInternal(err error) bool {
	ies := []error{
		ErrInvalidArgument,
		ErrUnauthenticated,
		ErrForbidden,
		ErrNotFound,
		ErrAlreadyExists,
		ErrFailedPrecondition,
		ErrResourceExhausted,
		ErrNotImplemented,
		ErrInternal,
		ErrCanceled,
		ErrUnavailable,
		ErrDeadlineExceeded,
		ErrOutOfRange,
	}
	_, ok := lo.Find(ies, func(ie error) bool {
		return errors.Is(err, ie)
	})
	return ok
}

func validationError(err error) error {
	const prefix = "validation"
	if err, ok := err.(validator.ValidationErrors); ok {
		return wrapError(prefix, ErrInvalidArgument, err)
	}
	return nil
}

func dbError(err error) error {
	const prefix = "database"

	//nolint:gocritic
	switch err := err.(type) {
	case *mysql.MySQLError:
		if err.Number == 1062 {
			return wrapError(prefix, ErrAlreadyExists, err)
		}
		return wrapError(prefix, ErrInternal, err)
	}

	switch {
	case errors.Is(err, gorm.ErrEmptySlice),
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrInvalidTransaction),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength),
		errors.Is(err, gorm.ErrMissingWhereClause),
		errors.Is(err, gorm.ErrModelValueRequired),
		errors.Is(err, gorm.ErrPrimaryKeyRequired):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, gorm.ErrNotImplemented):
		return wrapError(prefix, ErrNotImplemented, err)
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return wrapError(prefix, ErrInternal, err)
	default:
		return nil
	}
}

func authError(err error) error {
	const prefix = "auth"
	switch {
	case errors.Is(err, cognito.ErrInvalidArgument):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, cognito.ErrUnauthenticated):
		return wrapError(prefix, ErrUnauthenticated, err)
	case errors.Is(err, cognito.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, cognito.ErrAlreadyExists):
		return wrapError(prefix, ErrAlreadyExists, err)
	case errors.Is(err, cognito.ErrResourceExhausted):
		return wrapError(prefix, ErrResourceExhausted, err)
	case errors.Is(err, cognito.ErrCanceled):
		return wrapError(prefix, ErrCanceled, err)
	case errors.Is(err, cognito.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func storageError(err error) error {
	const prefix = "storage"
	switch {
	case errors.Is(err, storage.ErrInvalidURL):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, storage.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	default:
		return nil
	}
}

func mailerError(err error) error {
	const prefix = "mailer"
	switch {
	case errors.Is(err, mailer.ErrInvalidArgument):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, mailer.ErrUnauthenticated):
		return wrapError(prefix, ErrUnauthenticated, err)
	case errors.Is(err, mailer.ErrPermissionDenied):
		return wrapError(prefix, ErrForbidden, err)
	case errors.Is(err, mailer.ErrPayloadTooLong):
		return wrapError(prefix, ErrResourceExhausted, err)
	case errors.Is(err, mailer.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, mailer.ErrInternal):
		return wrapError(prefix, ErrInternal, err)
	case errors.Is(err, mailer.ErrUnavailable):
		return wrapError(prefix, ErrUnavailable, err)
	case errors.Is(err, mailer.ErrCanceled):
		return wrapError(prefix, ErrCanceled, err)
	case errors.Is(err, mailer.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func messagingError(err error) error {
	const prefix = "messaging"
	switch {
	case errors.Is(err, messaging.ErrInvalidArgument):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, messaging.ErrUnauthenticated):
		return wrapError(prefix, ErrUnauthenticated, err)
	case errors.Is(err, messaging.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, messaging.ErrResourceExhausted):
		return wrapError(prefix, ErrResourceExhausted, err)
	case errors.Is(err, messaging.ErrInternal):
		return wrapError(prefix, ErrInternal, err)
	case errors.Is(err, messaging.ErrCanceled):
		return wrapError(prefix, ErrCanceled, err)
	case errors.Is(err, messaging.ErrUnavailable):
		return wrapError(prefix, ErrUnavailable, err)
	case errors.Is(err, messaging.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func notifierError(err error) error {
	const prefix = "notifier"
	switch {
	case errors.Is(err, line.ErrInvalidArgument):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, line.ErrUnauthenticated):
		return wrapError(prefix, ErrUnauthenticated, err)
	case errors.Is(err, line.ErrPermissionDenied):
		return wrapError(prefix, ErrForbidden, err)
	case errors.Is(err, line.ErrPayloadTooLong):
		return wrapError(prefix, ErrResourceExhausted, err)
	case errors.Is(err, line.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, line.ErrAlreadyExists):
		return wrapError(prefix, ErrAlreadyExists, err)
	case errors.Is(err, line.ErrInternal):
		return wrapError(prefix, ErrInternal, err)
	case errors.Is(err, line.ErrUnavailable):
		return wrapError(prefix, ErrUnavailable, err)
	case errors.Is(err, line.ErrResourceExhausted):
		return wrapError(prefix, ErrResourceExhausted, err)
	case errors.Is(err, line.ErrCanceled):
		return wrapError(prefix, ErrCanceled, err)
	case errors.Is(err, line.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}

func externalError(err error) error {
	const prefix = "external"
	switch {
	case errors.Is(err, postalcode.ErrInvalidArgument):
		return wrapError(prefix, ErrInvalidArgument, err)
	case errors.Is(err, postalcode.ErrNotFound):
		return wrapError(prefix, ErrNotFound, err)
	case errors.Is(err, postalcode.ErrInternal):
		return wrapError(prefix, ErrInternal, err)
	case errors.Is(err, postalcode.ErrUnavailable):
		return wrapError(prefix, ErrUnavailable, err)
	case errors.Is(err, postalcode.ErrCanceled):
		return wrapError(prefix, ErrCanceled, err)
	case errors.Is(err, postalcode.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}
