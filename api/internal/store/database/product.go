package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const productTable = "products"

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

	stmt := p.db.Statement(ctx, p.db.DB, productTable, fields...)
	stmt = params.stmt(stmt)
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

func (p *product) Count(ctx context.Context, params *ListProductsParams) (int64, error) {
	total, err := p.db.Count(ctx, p.db.DB, &entity.Product{}, params.stmt)
	return total, exception.InternalError(err)
}

func (p *product) MultiGet(ctx context.Context, productIDs []string, fields ...string) (entity.Products, error) {
	var products entity.Products

	stmt := p.db.Statement(ctx, p.db.DB, productTable, fields...).
		Where("id IN (?)", productIDs)

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
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := product.FillJSON(); err != nil {
			return err
		}

		now := p.now()
		product.CreatedAt, product.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(productTable).Create(&product).Error
		return err
	})
	return exception.InternalError(err)
}

func (p *product) Update(ctx context.Context, productID string, params *UpdateProductParams) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := p.get(ctx, tx, productID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"producer_id":         params.ProducerID,
			"product_type_id":     params.TypeID,
			"name":                params.Name,
			"description":         params.Description,
			"public":              params.Public,
			"inventory":           params.Inventory,
			"weight":              params.Weight,
			"weight_unit":         params.WeightUnit,
			"item":                params.Item,
			"item_unit":           params.ItemUnit,
			"item_description":    params.ItemDescription,
			"price":               params.Price,
			"cost":                params.Cost,
			"expiration_date":     params.ExpirationDate,
			"storage_method_type": params.StorageMethodType,
			"delivery_type":       params.DeliveryType,
			"box60_rate":          params.Box60Rate,
			"box80_rate":          params.Box80Rate,
			"box100_rate":         params.Box100Rate,
			"origin_prefecture":   params.OriginPrefecture,
			"origin_city":         params.OriginCity,
			"updated_at":          p.now(),
		}
		if len(params.Media) > 0 {
			media, err := params.Media.Marshal()
			if err != nil {
				return fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["media"] = media
		}
		if len(params.TagIDs) > 0 {
			tagIDs, err := entity.ProductMarshalTagIDs(params.TagIDs)
			if err != nil {
				return fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["product_tag_id"] = tagIDs
		}
		if len(params.RecommendedPoints) > 0 {
			points, err := entity.ProductMarshalRecommendedPoints(params.RecommendedPoints)
			if err != nil {
				return fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["recommemded_points"] = points
		}
		err := tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", productID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (p *product) UpdateMedia(
	ctx context.Context, productID string, set func(media entity.MultiProductMedia) bool,
) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		product, err := p.get(ctx, tx, productID, "media")
		if err != nil {
			return err
		}
		if exists := set(product.Media); !exists {
			return fmt.Errorf("database: media is non-existent: %w", exception.ErrFailedPrecondition)
		}

		buf, err := product.Media.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"media":      datatypes.JSON(buf),
			"updated_at": p.now(),
		}

		err = tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", productID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (p *product) Delete(ctx context.Context, productID string) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := p.get(ctx, tx, productID); err != nil {
			return err
		}

		params := map[string]interface{}{
			"deleted_at": p.now(),
		}
		err := tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", productID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (p *product) get(ctx context.Context, tx *gorm.DB, productID string, fields ...string) (*entity.Product, error) {
	var product *entity.Product

	err := p.db.Statement(ctx, tx, productTable, fields...).
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
