package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const productTypeTable = "product_types"

type productType struct {
	db  *mysql.Client
	now func() time.Time
}

func newProductType(db *mysql.Client) database.ProductType {
	return &productType{
		db:  db,
		now: jst.Now,
	}
}

type listProductTypesParams database.ListProductTypesParams

func (p listProductTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", p.CategoryID)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	if len(p.Orders) == 0 {
		stmt = stmt.Order("category_id ASC, name ASC")
	}
	return stmt
}

func (p listProductTypesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *productType) List(
	ctx context.Context, params *database.ListProductTypesParams, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes

	p := listProductTypesParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, productTypeTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&productTypes).Error; err != nil {
		return nil, dbError(err)
	}
	if err := productTypes.Fill(); err != nil {
		return nil, dbError(err)
	}
	return productTypes, nil
}

func (t *productType) Count(ctx context.Context, params *database.ListProductTypesParams) (int64, error) {
	p := listProductTypesParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.ProductType{}, p.stmt)
	return total, dbError(err)
}

func (t *productType) MultiGet(
	ctx context.Context, productTypeIDs []string, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes

	stmt := t.db.Statement(ctx, t.db.DB, productTypeTable, fields...).
		Where("id IN (?)", productTypeIDs)

	if err := stmt.Find(&productTypes).Error; err != nil {
		return nil, dbError(err)
	}
	if err := productTypes.Fill(); err != nil {
		return nil, dbError(err)
	}
	return productTypes, nil
}

func (t *productType) Get(ctx context.Context, productTypeID string, fields ...string) (*entity.ProductType, error) {
	productType, err := t.get(ctx, t.db.DB, productTypeID, fields...)
	return productType, dbError(err)
}

func (t *productType) Create(ctx context.Context, productType *entity.ProductType) error {
	now := t.now()
	productType.CreatedAt, productType.UpdatedAt = now, now

	err := t.db.DB.WithContext(ctx).Table(productTypeTable).Create(&productType).Error
	return dbError(err)
}

func (t *productType) Update(ctx context.Context, productTypeID, name, iconURL string) error {
	params := map[string]interface{}{
		"name":       name,
		"icon_url":   iconURL,
		"updated_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).
		Table(productTypeTable).
		Where("id = ?", productTypeID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (t *productType) UpdateIcons(ctx context.Context, productTypeID string, icons common.Images) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		productType, err := t.get(ctx, tx, productTypeID, "icon_url")
		if err != nil {
			return err
		}
		if productType.IconURL == "" {
			return fmt.Errorf("database: icon url is empty: %w", database.ErrFailedPrecondition)
		}

		buf, err := icons.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"icons":      datatypes.JSON(buf),
			"updated_at": t.now(),
		}

		err = tx.WithContext(ctx).
			Table(productTypeTable).
			Where("id = ?", productTypeID).
			Updates(params).Error
		return err
	})
	return dbError(err)
}

func (t *productType) Delete(ctx context.Context, productTypeID string) error {
	stmt := t.db.DB.WithContext(ctx).
		Table(productTypeTable).
		Where("id = ?", productTypeID)

	err := stmt.Delete(&entity.ProductType{}).Error
	return dbError(err)
}

func (t *productType) get(
	ctx context.Context, tx *gorm.DB, productTypeID string, fields ...string,
) (*entity.ProductType, error) {
	var productType *entity.ProductType

	stmt := t.db.Statement(ctx, tx, productTypeTable, fields...).
		Where("id = ?", productTypeID)

	if err := stmt.First(&productType).Error; err != nil {
		return nil, err
	}
	if err := productType.Fill(); err != nil {
		return nil, err
	}
	return productType, nil
}
