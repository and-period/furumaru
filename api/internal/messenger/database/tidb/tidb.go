package tidb

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/pkg/mysql"
	gmysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(db *mysql.Client) *database.Database {
	return &database.Database{
		Contact:         NewContact(db),
		ContactCategory: NewContactCategory(db),
		ContactRead:     NewContactRead(db),
		Message:         NewMessage(db),
		MessageTemplate: NewMessageTemplate(db),
		Notification:    NewNotification(db),
		PushTemplate:    NewPushTemplate(db),
		ReceivedQueue:   NewReceivedQueue(db),
		ReportTemplate:  NewReportTemplate(db),
		Schedule:        NewSchedule(db),
	}
}

func dbError(err error) error {
	var derr *database.Error
	if err == nil || errors.As(err, &derr) {
		return err
	}

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", database.ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", database.ErrDeadlineExceeded, err.Error())
	}

	//nolint:errorlint,gocritic
	switch err := err.(type) {
	case *gmysql.MySQLError:
		if err.Number == 1062 {
			return fmt.Errorf("%w: %s", database.ErrAlreadyExists, err.Error())
		}
		return fmt.Errorf("%w: %s", database.ErrInternal, err.Error())
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
		return fmt.Errorf("%w: %s", database.ErrInvalidArgument, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%w: %s", database.ErrNotFound, err.Error())
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return fmt.Errorf("%w: %s", database.ErrInternal, err.Error())
	default:
		return fmt.Errorf("%w: %s", database.ErrUnknown, err.Error())
	}
}
