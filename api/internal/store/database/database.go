//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/database"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Staff Staff
	Store Store
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Staff: NewStaff(params.Database),
		Store: NewStore(params.Database),
	}
}

/**
 * interface
 */
type Staff interface {
	ListByStoreID(ctx context.Context, storeID int64, fields ...string) (entity.Staffs, error)
}

type Store interface {
	List(ctx context.Context, params *ListStoresParams, fields ...string) (entity.Stores, error)
	Get(ctx context.Context, storeID int64, fields ...string) (*entity.Store, error)
	Create(ctx context.Context, store *entity.Store) error
	Update(ctx context.Context, storeID int64, name, thumbnailURL string) error
}

/**
 * params
 */
type ListStoresParams struct {
	Limit  int
	Offset int
}
