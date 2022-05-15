//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/marche/api/internal/store/entity"
)

//nolint:revive
type StoreService interface {
	// 店舗一覧取得
	ListStores(ctx context.Context, in *ListStoresInput) (entity.Stores, error)
	// 店舗取得
	GetStore(ctx context.Context, in *GetStoreInput) (*entity.Store, error)
	// 店舗登録
	CreateStore(ctx context.Context, in *CreateStoreInput) (*entity.Store, error)
	// 店舗更新
	UpdateStore(ctx context.Context, in *UpdateStoreInput) error
	// 店舗のスタッフ一覧取得
	ListStaffsByStoreID(ctx context.Context, in *ListStaffsByStoreIDInput) (entity.Staffs, error)
}
