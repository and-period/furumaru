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
	Product     Product
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

type Product interface {
	List(ctx context.Context, params *ListProductsParams, fields ...string) (entity.Products, error)
	Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, productID string, params *UpdateProductParams) error
	Delete(ctx context.Context, productID string) error
}

type ProductType interface {
	List(ctx context.Context, params *ListProductTypesParams, fields ...string) (entity.ProductTypes, error)
	MultiGet(ctx context.Context, productTypeIDs []string, fields ...string) (entity.ProductTypes, error)
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

type ListProductsParams struct {
	Name       string
	ProducerID string
	CreatedBy  string
	Limit      int
	Offset     int
}

type UpdateProductParams struct {
	ProducerID       string
	CategoryID       string
	TypeID           string
	Name             string
	Description      string
	Public           bool
	Inventory        int64
	Weight           int64
	WeightUnit       entity.WeightUnit
	Item             int64
	ItemUnit         string
	ItemDescription  string
	Media            entity.MultiProductMedia
	Price            int64
	DeliveryType     entity.DeliveryType
	Box60Rate        int64
	Box80Rate        int64
	Box100Rate       int64
	OriginPrefecture string
	OriginCity       string
	UpdatedBy        string
}

type ListProductTypesParams struct {
	Name       string
	CategoryID string
	Limit      int
	Offset     int
}
