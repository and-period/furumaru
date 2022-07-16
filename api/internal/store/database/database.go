//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"gorm.io/gorm"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Category    Category
	Product     Product
	ProductType ProductType
	Shipping    Shipping
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Category:    NewCategory(params.Database),
		Product:     NewProduct(params.Database),
		ProductType: NewProductType(params.Database),
		Shipping:    NewShipping(params.Database),
	}
}

/**
 * interface
 */
type Category interface {
	List(ctx context.Context, params *ListCategoriesParams, fields ...string) (entity.Categories, error)
	Count(ctx context.Context, params *ListCategoriesParams) (int64, error)
	MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.Categories, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, categoryID, name string) error
	Delete(ctx context.Context, categoryID string) error
}

type Product interface {
	List(ctx context.Context, params *ListProductsParams, fields ...string) (entity.Products, error)
	Count(ctx context.Context, params *ListProductsParams) (int64, error)
	Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, productID string, params *UpdateProductParams) error
	Delete(ctx context.Context, productID string) error
}

type ProductType interface {
	List(ctx context.Context, params *ListProductTypesParams, fields ...string) (entity.ProductTypes, error)
	Count(ctx context.Context, params *ListProductTypesParams) (int64, error)
	MultiGet(ctx context.Context, productTypeIDs []string, fields ...string) (entity.ProductTypes, error)
	Create(ctx context.Context, productType *entity.ProductType) error
	Update(ctx context.Context, productTypeID, name string) error
	Delete(ctx context.Context, productTypeID string) error
}

type Shipping interface {
	List(ctx context.Context, params *ListShippingsParams, fields ...string) (entity.Shippings, error)
	Count(ctx context.Context, params *ListShippingsParams) (int64, error)
	Get(ctx context.Context, shoppingID string, fields ...string) (*entity.Shipping, error)
	Create(ctx context.Context, shipping *entity.Shipping) error
	Update(ctx context.Context, shippingID string, params *UpdateShippingParams) error
	Delete(ctx context.Context, shippingID string) error
}

/**
 * params
 */
type ListCategoriesParams struct {
	Name   string
	Limit  int
	Offset int
}

func (p *ListCategoriesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	return stmt
}

type ListShippingsParams struct {
	Limit  int
	Offset int
}

type UpdateShippingParams struct {
	Name               string
	Box60Rates         entity.ShippingRates
	Box60Refrigerated  int64
	Box60Frozen        int64
	Box80Rates         entity.ShippingRates
	Box80Refrigerated  int64
	Box80Frozen        int64
	Box100Rates        entity.ShippingRates
	Box100Refrigerated int64
	Box100Frozen       int64
	HasFreeShipping    bool
	FreeShippingRates  int64
}

type ListProductsParams struct {
	Name       string
	ProducerID string
	CreatedBy  string
	Limit      int
	Offset     int
}

func (p *ListProductsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", p.ProducerID)
	}
	if p.CreatedBy != "" {
		stmt = stmt.Where("created_by = ?", p.CreatedBy)
	}
	return stmt
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

func (p *ListProductTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", p.CategoryID)
	}
	return stmt
}
