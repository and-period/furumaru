package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const productTypeTable = "product_types"

type productType struct {
	db  *database.Client
	now func() time.Time
}

func NewProductType(db *database.Client) ProductType {
	return &productType{
		db:  db,
		now: jst.Now,
	}
}

func (t *productType) List(
	ctx context.Context, params *ListProductTypesParams, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes

	stmt := t.db.Statement(ctx, t.db.DB, productTypeTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&productTypes).Error
	return productTypes, exception.InternalError(err)
}

func (t *productType) Count(ctx context.Context, params *ListProductTypesParams) (int64, error) {
	var total int64

	stmt := t.db.Count(ctx, t.db.DB, productTypeTable)
	stmt = params.stmt(stmt)

	err := stmt.Find(&total).Error
	return total, exception.InternalError(err)
}

func (t *productType) MultiGet(
	ctx context.Context, productTypeIDs []string, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes

	err := t.db.Statement(ctx, t.db.DB, productTypeTable, fields...).
		Where("id IN (?)", productTypeIDs).
		Find(&productTypes).Error
	return productTypes, exception.InternalError(err)
}

func (t *productType) Get(ctx context.Context, productTypeID string, fields ...string) (*entity.ProductType, error) {
	productType, err := t.get(ctx, t.db.DB, productTypeID, fields...)
	return productType, exception.InternalError(err)
}

func (t *productType) Create(ctx context.Context, productType *entity.ProductType) error {
	_, err := t.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		err := tx.WithContext(ctx).
			Table(categoryTable).Select(categoryFields).
			Where("id = ?", productType.CategoryID).
			First(&entity.Category{}).Error
		if err != nil {
			return nil, err
		}

		now := t.now()
		productType.CreatedAt, productType.UpdatedAt = now, now

		err = tx.WithContext(ctx).Table(productTypeTable).Create(&productType).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (t *productType) Update(ctx context.Context, productTypeID, name, iconURL string) error {
	_, err := t.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := t.get(ctx, tx, productTypeID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"name":       name,
			"icon_url":   iconURL,
			"updated_at": t.now(),
		}
		err := tx.WithContext(ctx).
			Table(productTypeTable).
			Where("id = ?", productTypeID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (t *productType) Delete(ctx context.Context, productTypeID string) error {
	_, err := t.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := t.get(ctx, tx, productTypeID); err != nil {
			return nil, err
		}

		err := tx.WithContext(ctx).
			Table(productTypeTable).
			Where("id = ?", productTypeID).
			Delete(&entity.ProductType{}).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (t *productType) get(
	ctx context.Context, tx *gorm.DB, productTypeID string, fields ...string,
) (*entity.ProductType, error) {
	var productType *entity.ProductType

	err := t.db.Statement(ctx, tx, productTypeTable, fields...).
		Where("id = ?", productTypeID).
		First(&productType).Error
	return productType, err
}
