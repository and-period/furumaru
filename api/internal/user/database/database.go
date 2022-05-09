//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrInvalidArgument    = errors.New("database: invalid argument")
	ErrNotFound           = errors.New("database: not found")
	ErrFailedPrecondition = errors.New("database: failed precondition")
	ErrAlreadyExists      = errors.New("database: already exists")
	ErrNotImplemented     = errors.New("database: not implemented")
	ErrInternal           = errors.New("database: internal")
	ErrUnknown            = errors.New("database: unknown")
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Admin Admin
	Shop  Shop
	User  User
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Admin: NewAdmin(params.Database),
		Shop:  NewShop(params.Database),
		User:  NewUser(params.Database),
	}
}

/**
 * interface
 */
type Admin interface {
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	UpdateEmail(ctx context.Context, adminID, email string) error
}

type Shop interface {
	List(ctx context.Context, params *ListShopsParams, fields ...string) (entity.Shops, error)
	MultiGet(ctx context.Context, shopIDs []string, fields ...string) (entity.Shops, error)
	Get(ctx context.Context, shopID string, fields ...string) (*entity.Shop, error)
	Create(ctx context.Context, shop *entity.Shop) error
}

type User interface {
	Get(ctx context.Context, userID string, fields ...string) (*entity.User, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	UpdateVerified(ctx context.Context, userID string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	Delete(ctx context.Context, userID string) error
}

/**
 * params
 */
type ListShopsParams struct {
	Limit  int
	Offset int
}

/**
 * private methods
 */
func dbError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrFailedPrecondition) {
		return err
	}

	//nolint:gocritic
	switch err := err.(type) {
	case *mysql.MySQLError:
		if err.Number == 1062 {
			return fmt.Errorf("%w: %s", ErrAlreadyExists, err)
		}
		return fmt.Errorf("%w: %s", ErrUnknown, err)
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
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err)
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%w: %s", ErrNotFound, err)
	case errors.Is(err, gorm.ErrNotImplemented):
		return fmt.Errorf("%w: %s", ErrNotImplemented, err)
	case errors.Is(err, gorm.ErrDryRunModeUnsupported),
		errors.Is(err, gorm.ErrInvalidDB),
		errors.Is(err, gorm.ErrRegistered),
		errors.Is(err, gorm.ErrUnsupportedDriver),
		errors.Is(err, gorm.ErrUnsupportedRelation):
		return fmt.Errorf("%w: %s", ErrInternal, err)
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err)
	}
}
