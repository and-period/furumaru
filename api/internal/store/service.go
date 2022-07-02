//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package store

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Service interface {
	// カテゴリ一覧取得
	ListCategories(ctx context.Context, in *ListCategoriesInput) (entity.Categories, error)
	// カテゴリ一覧取得(ID指定)
	MultiGetCategories(ctx context.Context, in *MultiGetCategoriesInput) (entity.Categories, error)
	// カテゴリ登録
	CreateCategory(ctx context.Context, in *CreateCategoryInput) (*entity.Category, error)
	// カテゴリ更新
	UpdateCategory(ctx context.Context, in *UpdateCategoryInput) error
	// カテゴリ削除
	DeleteCategory(ctx context.Context, in *DeleteCategoryInput) error
	// 品目一覧取得
	ListProductTypes(ctx context.Context, in *ListProductTypesInput) (entity.ProductTypes, error)
	// 品目一覧取得(ID指定)
	MultiGetProductTypes(ctx context.Context, in *MultiGetProductTypesInput) (entity.ProductTypes, error)
	// 品目登録
	CreateProductType(ctx context.Context, in *CreateProductTypeInput) (*entity.ProductType, error)
	// 品目更新
	UpdateProductType(ctx context.Context, in *UpdateProductTypeInput) error
	// 品目削除
	DeleteProductType(ctx context.Context, in *DeleteProductTypeInput) error
	// 商品一覧取得
	ListProducts(ctx context.Context, in *ListProductsInput) (entity.Products, error)
	// 商品取得
	GetProduct(ctx context.Context, in *GetProductInput) (*entity.Product, error)
	// 商品登録
	CreateProduct(ctx context.Context, in *CreateProductInput) (*entity.Product, error)
	// 商品更新
	UpdateProduct(ctx context.Context, in *UpdateProductInput) error
	// 商品削除
	DeleteProduct(ctx context.Context, in *DeleteProductInput) error
}
