package exception

import (
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/line"
	"github.com/and-period/furumaru/api/pkg/mailer"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
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
	if err := notifierError(err); err != nil {
		return err
	}
	return wrapError("internal", ErrUnknown, err)
}

func Retryable(err error) bool {
	return errors.Is(err, ErrCanceled) ||
		errors.Is(err, ErrUnavailable) ||
		errors.Is(err, ErrDeadlineExceeded) ||
		errors.Is(err, ErrOutOfRange)
}

func wrapError(prefix string, code, err error) error {
	return fmt.Errorf("%s: %w: %s", prefix, code, err.Error())
}

func isInternal(err error) bool {
	return errors.Is(err, ErrInvalidArgument) ||
		errors.Is(err, ErrUnauthenticated) ||
		errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrAlreadyExists) ||
		errors.Is(err, ErrFailedPrecondition) ||
		errors.Is(err, ErrResourceExhausted) ||
		errors.Is(err, ErrNotImplemented) ||
		errors.Is(err, ErrInternal)
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
	case errors.Is(err, mailer.ErrTimeout):
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
	case errors.Is(err, line.ErrTimeout):
		return wrapError(prefix, ErrDeadlineExceeded, err)
	default:
		return nil
	}
}
