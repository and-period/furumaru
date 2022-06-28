package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const productTable = "products"

var productFields = []string{
	"id", "producer_id", "category_id", "product_type_id",
	"name", "description", "public", "inventory",
	"weight", "weight_unit", "item", "item_unit", "item_description",
	"media", "price", "delivery_type", "box60_rate", "box80_rate", "box100_rate",
	"origin_prefecture", "origin_city", "created_by", "updated_by",
	"created_at", "updated_at", "deleted_at",
}

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
	var products entity.Products
	if len(fields) == 0 {
		fields = productFields
	}

	stmt := p.db.DB.WithContext(ctx).Table(productTable).Select(fields)
	if params.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", params.Name))
	}
	if params.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", params.ProducerID)
	}
	if params.CreatedBy != "" {
		stmt = stmt.Where("created_by = ?", params.CreatedBy)
	}
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&products).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := products.Fill(); err != nil {
		return nil, exception.InternalError(err)
	}
	return products, nil
}

func (p *product) Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error) {
	product, err := p.get(ctx, p.db.DB, productID, fields...)
	return product, exception.InternalError(err)
}

func (p *product) Create(ctx context.Context, product *entity.Product) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if err := product.FillJSON(); err != nil {
			return nil, err
		}

		now := p.now()
		product.CreatedAt, product.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(productTable).Create(&product).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *product) Update(ctx context.Context, product *entity.Product) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, product.ID); err != nil {
			return nil, err
		}
		if err := product.FillJSON(); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"producer_id":       product.ProducerID,
			"category_id":       product.CategoryID,
			"product_type_id":   product.TypeID,
			"name":              product.Name,
			"description":       product.Description,
			"public":            product.Public,
			"inventory":         product.Inventory,
			"weight":            product.Weight,
			"weight_unit":       product.WeightUnit,
			"item":              product.Item,
			"item_unit":         product.ItemUnit,
			"item_description":  product.ItemDescription,
			"media":             product.MediaJSON,
			"price":             product.Price,
			"delivery_type":     product.DeliveryType,
			"box60_rate":        product.Box60Rate,
			"box80_rate":        product.Box80Rate,
			"box100_rate":       product.Box100Rate,
			"origin_prefecture": product.OriginPrefecture,
			"origin_city":       product.OriginCity,
			"updated_at":        p.now(),
			"updated_by":        product.UpdatedBy,
		}
		err := tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", product.ID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *product) Delete(ctx context.Context, productID string) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, productID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"deleted_at": p.now(),
		}
		err := tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", productID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *product) get(ctx context.Context, tx *gorm.DB, productID string, fields ...string) (*entity.Product, error) {
	var product *entity.Product
	if len(fields) == 0 {
		fields = productFields
	}

	err := tx.WithContext(ctx).
		Table(productTable).Select(fields).
		Where("id = ?", productID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	if err := product.Fill(); err != nil {
		return nil, err
	}
	return product, nil
}
