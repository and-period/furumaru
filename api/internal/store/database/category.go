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

const categoryTable = "categories"

var categoryFields = []string{
	"id", "name", "created_at", "updated_at",
}

type category struct {
	db  *database.Client
	now func() time.Time
}

func NewCategory(db *database.Client) Category {
	return &category{
		db:  db,
		now: jst.Now,
	}
}

func (c *category) List(
	ctx context.Context, params *ListCategoriesParams, fields ...string,
) (entity.Categories, error) {
	var categories entity.Categories
	if len(fields) == 0 {
		fields = categoryFields
	}

	stmt := c.db.DB.WithContext(ctx).Table(categoryTable).Select(fields)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&categories).Error
	return categories, exception.InternalError(err)
}

func (c *category) Count(ctx context.Context, params *ListCategoriesParams) (int64, error) {
	var total int64

	stmt := c.db.DB.WithContext(ctx).Table(categoryTable).Select("COUNT(*)")
	stmt = params.stmt(stmt)

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (c *category) MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.Categories, error) {
	var categories entity.Categories
	if len(fields) == 0 {
		fields = categoryFields
	}

	err := c.db.DB.WithContext(ctx).
		Table(categoryTable).Select(fields).
		Where("id IN (?)", categoryIDs).
		Find(&categories).Error
	return categories, exception.InternalError(err)
}

func (c *category) Create(ctx context.Context, category *entity.Category) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := c.now()
		category.CreatedAt, category.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(categoryTable).Create(&category).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *category) Update(ctx context.Context, categoryID, name string) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, categoryID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"name":       name,
			"updated_at": c.now(),
		}
		err := tx.WithContext(ctx).
			Table(categoryTable).
			Where("id = ?", categoryID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *category) Delete(ctx context.Context, categoryID string) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, categoryID); err != nil {
			return nil, err
		}

		err := tx.WithContext(ctx).
			Table(categoryTable).
			Where("id = ?", categoryID).
			Delete(&entity.Category{}).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *category) get(
	ctx context.Context, tx *gorm.DB, categoryID string, fields ...string,
) (*entity.Category, error) {
	var category *entity.Category
	if len(fields) == 0 {
		fields = categoryFields
	}

	err := tx.WithContext(ctx).
		Table(categoryTable).Select(fields).
		Where("id = ?", categoryID).
		First(&category).Error
	return category, err
}
