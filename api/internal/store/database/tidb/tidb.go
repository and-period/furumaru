package tidb

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/database/mysql"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	gmysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(db *apmysql.Client) *database.Database {
	client := mysql.NewDatabase(db)
	return &database.Database{
		CartActionLog:            client.CartActionLog,
		Category:                 NewCategory(db, client.Category),
		Experience:               NewExperience(db),
		ExperienceReview:         client.ExperienceReview,
		ExperienceReviewReaction: client.ExperienceReviewReaction,
		ExperienceType:           NewExperienceType(db, client.ExperienceType),
		Live:                     client.Live,
		Order:                    NewOrder(db, client.Order),
		PaymentSystem:            client.PaymentSystem,
		Product:                  NewProduct(db, client.Product),
		ProductReview:            client.ProductReview,
		ProductReviewReaction:    client.ProductReviewReaction,
		ProductTag:               NewProductTag(db, client.ProductTag),
		ProductType:              NewProductType(db, client.ProductType),
		Promotion:                NewPromotion(db, client.Promotion),
		Schedule:                 client.Schedule,
		Shipping:                 client.Shipping,
		Spot:                     NewSpot(db, client.Spot),
		SpotType:                 NewSpotType(db, client.SpotType),
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

	//nolint:gocritic
	switch err := err.(type) {
	case *gmysql.MySQLError:
		if err.Number == 1062 {
			return fmt.Errorf("%w: %s", database.ErrAlreadyExists, err)
		}
		return fmt.Errorf("%w: %s", database.ErrInternal, err)
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
		return fmt.Errorf("%w: %s", database.ErrInvalidArgument, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%w: %s", database.ErrNotFound, err)
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return fmt.Errorf("%w: %s", database.ErrInternal, err)
	default:
		return fmt.Errorf("%w: %s", database.ErrUnknown, err)
	}
}
