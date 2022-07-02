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

const productTypeTable = "product_types"

var productTypeFields = []string{
	"id", "name", "category_id", "created_at", "updated_at",
}

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
	if len(fields) == 0 {
		fields = productTypeFields
	}

	stmt := t.db.DB.WithContext(ctx).Table(productTypeTable).Select(fields)
	if params.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", params.Name))
	}
	if params.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", params.CategoryID)
	}
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&productTypes).Error
	return productTypes, exception.InternalError(err)
}

func (t *productType) MultiGet(
	ctx context.Context, productTypeIDs []string, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes
	if len(fields) == 0 {
		fields = productTypeFields
	}

	err := t.db.DB.WithContext(ctx).
		Table(productTypeTable).Select(fields).
		Where("id IN (?)", productTypeIDs).
		Find(&productTypes).Error
	return productTypes, exception.InternalError(err)
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

func (t *productType) Update(ctx context.Context, productTypeID, name string) error {
	_, err := t.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := t.get(ctx, tx, productTypeID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"name":       name,
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
	if len(fields) == 0 {
		fields = productTypeFields
	}

	err := tx.WithContext(ctx).
		Table(productTypeTable).Select(fields).
		Where("id = ?", productTypeID).
		First(&productType).Error
	return productType, err
}
