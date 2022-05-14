//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/marche/api/internal/store/entity"
)

//nolint:revive
type StoreService interface {
	ListStaffsByStoreID(ctx context.Context, in *ListStaffsByStoreIDInput) (entity.Staffs, error)
	ListStores(ctx context.Context, in *ListStoresInput) (entity.Stores, error)
	GetStore(ctx context.Context, in *GetStoreInput) (*entity.Store, error)
	CreateStore(ctx context.Context, in *CreateStoreInput) (*entity.Store, error)
	UpdateStore(ctx context.Context, in *UpdateStoreInput) error
}
