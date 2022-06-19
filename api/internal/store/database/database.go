//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Category    Category
	ProductType ProductType
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Category:    NewCategory(params.Database),
		ProductType: NewProductType(params.Database),
	}
}

/**
 * interface
 */
type Category interface {
	List(ctx context.Context, params *ListCategoriesParams, fields ...string) (entity.Categories, error)
	MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.Categories, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, categoryID, name string) error
	Delete(ctx context.Context, categoryID string) error
}

type ProductType interface {
	List(ctx context.Context, params *ListProductTypesParams, fields ...string) (entity.ProductTypes, error)
	Create(ctx context.Context, productType *entity.ProductType) error
	Update(ctx context.Context, productTypeID, name string) error
	Delete(ctx context.Context, productTypeID string) error
}

/**
 * params
 */
type ListCategoriesParams struct {
	Name   string
	Limit  int
	Offset int
}

type ListProductTypesParams struct {
	Name       string
	CategoryID string
	Limit      int
	Offset     int
}
