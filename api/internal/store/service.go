//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Service interface {
	// カテゴリ一覧取得
	ListCategories(ctx context.Context, in *ListCategoriesInput) (entity.Categories, int64, error)
	// カテゴリ一覧取得(ID指定)
	MultiGetCategories(ctx context.Context, in *MultiGetCategoriesInput) (entity.Categories, error)
	// カテゴリ取得
	GetCategory(ctx context.Context, in *GetCategoryInput) (*entity.Category, error)
	// カテゴリ登録
	CreateCategory(ctx context.Context, in *CreateCategoryInput) (*entity.Category, error)
	// カテゴリ更新
	UpdateCategory(ctx context.Context, in *UpdateCategoryInput) error
	// カテゴリ削除
	DeleteCategory(ctx context.Context, in *DeleteCategoryInput) error
	// 品目一覧取得
	ListProductTypes(ctx context.Context, in *ListProductTypesInput) (entity.ProductTypes, int64, error)
	// 品目一覧取得(ID指定)
	MultiGetProductTypes(ctx context.Context, in *MultiGetProductTypesInput) (entity.ProductTypes, error)
	// 品目取得
	GetProductType(ctx context.Context, in *GetProductTypeInput) (*entity.ProductType, error)
	// 品目登録
	CreateProductType(ctx context.Context, in *CreateProductTypeInput) (*entity.ProductType, error)
	// 品目更新
	UpdateProductType(ctx context.Context, in *UpdateProductTypeInput) error
	// 品目削除
	DeleteProductType(ctx context.Context, in *DeleteProductTypeInput) error
	// 配送設定一覧取得
	ListShippings(ctx context.Context, in *ListShippingsInput) (entity.Shippings, int64, error)
	// 配送設定取得
	GetShipping(ctx context.Context, in *GetShippingInput) (*entity.Shipping, error)
	// 配送設定登録
	CreateShipping(ctx context.Context, in *CreateShippingInput) (*entity.Shipping, error)
	// 配送設定更新
	UpdateShipping(ctx context.Context, in *UpdateShippingInput) error
	// 配送設定削除
	DeleteShipping(ctx context.Context, in *DeleteShippingInput) error
	// 商品一覧取得
	ListProducts(ctx context.Context, in *ListProductsInput) (entity.Products, int64, error)
	// 商品取得
	GetProduct(ctx context.Context, in *GetProductInput) (*entity.Product, error)
	// 商品登録
	CreateProduct(ctx context.Context, in *CreateProductInput) (*entity.Product, error)
	// 商品更新
	UpdateProduct(ctx context.Context, in *UpdateProductInput) error
	// 商品削除
	DeleteProduct(ctx context.Context, in *DeleteProductInput) error
	// プロモーション一覧取得
	ListPromotions(ctx context.Context, in *ListPromotionsInput) (entity.Promotions, int64, error)
	// プロモーション取得
	GetPromotion(ctx context.Context, in *GetPromotionInput) (*entity.Promotion, error)
	// プロモーション登録
	CreatePromotion(ctx context.Context, in *CreatePromotionInput) (*entity.Promotion, error)
	// プロモーション更新
	UpdatePromotion(ctx context.Context, in *UpdatePromotionInput) error
	// プロモーション削除
	DeletePromotion(ctx context.Context, in *DeletePromotionInput) error
	// 開催スケジュール登録
	CreateSchedule(ctx context.Context, in *CreateScheduleInput) (*entity.Schedule, entity.Lives, error)
	// 注文履歴一覧取得
	ListOrders(ctx context.Context, in *ListOrdersInput) (entity.Orders, int64, error)
	// 注文履歴取得
	GetOrder(ctx context.Context, in *GetOrderInput) (*entity.Order, error)
}
