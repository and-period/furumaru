//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Admin Admin
	User  User
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Admin: NewAdmin(params.Database),
		User:  NewUser(params.Database),
	}
}

/**
 * interface
 */
type Admin interface {
	List(ctx context.Context, params *ListAdminsParams, fields ...string) (entity.Admins, error)
	MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.Admins, error)
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	Create(ctx context.Context, admin *entity.Admin) error
	UpdateEmail(ctx context.Context, adminID, email string) error
}

type User interface {
	Get(ctx context.Context, userID string, fields ...string) (*entity.User, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	UpdateVerified(ctx context.Context, userID string) error
	UpdateAccount(ctx context.Context, userID, accountID, username string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	Delete(ctx context.Context, userID string) error
}

/**
 * params
 */
type ListAdminsParams struct {
	Roles  []entity.AdminRole
	Limit  int
	Offset int
}
