//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

//nolint:revive
type StoreService interface {
	// カテゴリ一覧取得
	ListCategories(ctx context.Context, in *ListCategoriesInput) (entity.Categories, error)
	// カテゴリ登録
	CreateCategory(ctx context.Context, in *CreateCategoryInput) (*entity.Category, error)
	// カテゴリ更新
	UpdateCategory(ctx context.Context, in *UpdateCategoryInput) error
	// カテゴリ削除
	DeleteCategory(ctx context.Context, in *DeleteCategoryInput) error
	// 品目一覧取得
	ListProductTypes(ctx context.Context, in *ListProductTypesInput) (entity.ProductTypes, error)
	// 品目登録
	CreateProductType(ctx context.Context, in *CreateProductTypeInput) (*entity.ProductType, error)
	// 品目更新
	UpdateProductType(ctx context.Context, in *UpdateProductTypeInput) error
	// 品目削除
	DeleteProductType(ctx context.Context, in *DeleteProductTypeInput) error
}
