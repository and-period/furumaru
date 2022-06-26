package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// const productTable = "products"

// var productFields = []string{
// 	"id", "producer_id", "category_id", "product_type_id",
// 	"name", "description", "public", "inventory",
// 	"weight", "weight_unit", "item", "item_unit", "item_description",
// 	"media", "price", "delivery_type", "box60_rate", "box80_rate", "box100_rate",
// 	"origin_prefecture", "origin_city", "created_by", "updated_by",
// 	"created_at", "updated_at", "deleted_at",
// }

type product struct {
	db  *database.Client
	now func() time.Time
}

func NewProduct(db *database.Client) Product {
	return &product{
		db:  db,
		now: jst.Now,
	}
}

func (p *product) List(ctx context.Context, params *ListProductsParams, fields ...string) (entity.Products, error) {
	return nil, exception.ErrNotImplemented
}

func (p *product) Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error) {
	return nil, exception.ErrNotImplemented
}

func (p *product) Create(ctx context.Context, product *entity.Product) error {
	return exception.ErrNotImplemented
}

func (p *product) Update(ctx context.Context, product *entity.Product) error {
	return exception.ErrNotImplemented
}

func (p *product) Delete(ctx context.Context, productID string) error {
	return exception.ErrNotImplemented
}
