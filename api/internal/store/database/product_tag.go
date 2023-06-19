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

const productTagTable = "product_tags"

type productTag struct {
	db  *database.Client
	now func() time.Time
}

func NewProductTag(db *database.Client) ProductTag {
	return &productTag{
		db:  db,
		now: jst.Now,
	}
}

func (t *productTag) List(
	ctx context.Context, params *ListProductTagsParams, fields ...string,
) (entity.ProductTags, error) {
	var tags entity.ProductTags

	stmt := t.db.Statement(ctx, t.db.DB, productTagTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&tags).Error
	return tags, exception.InternalError(err)
}

func (t *productTag) Count(ctx context.Context, params *ListProductTagsParams) (int64, error) {
	total, err := t.db.Count(ctx, t.db.DB, &entity.ProductTag{}, params.stmt)
	return total, exception.InternalError(err)
}

func (t *productTag) MultiGet(
	ctx context.Context, productTagIDs []string, fields ...string,
) (entity.ProductTags, error) {
	var tags entity.ProductTags

	err := t.db.Statement(ctx, t.db.DB, productTagTable, fields...).
		Where("id IN (?)", productTagIDs).
		Find(&tags).Error
	return tags, exception.InternalError(err)
}

func (t *productTag) Get(ctx context.Context, productTagID string, fields ...string) (*entity.ProductTag, error) {
	tag, err := t.get(ctx, t.db.DB, productTagID, fields...)
	return tag, exception.InternalError(err)
}

func (t *productTag) Create(ctx context.Context, tag *entity.ProductTag) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := t.now()
		tag.CreatedAt, tag.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(productTagTable).Create(&tag).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *productTag) Update(ctx context.Context, productTagID, name string) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := t.get(ctx, tx, productTagID); err != nil {
			return err
		}

		params := map[string]interface{}{
			"name":       name,
			"updated_at": t.now(),
		}
		err := tx.WithContext(ctx).
			Table(productTagTable).
			Where("id = ?", productTagID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *productTag) Delete(ctx context.Context, productTagID string) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := t.get(ctx, tx, productTagID); err != nil {
			return err
		}

		err := tx.WithContext(ctx).
			Table(productTagTable).
			Where("id = ?", productTagID).
			Delete(&entity.ProductTag{}).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *productTag) get(
	ctx context.Context, tx *gorm.DB, productTagID string, fields ...string,
) (*entity.ProductTag, error) {
	var tag *entity.ProductTag

	stmt := t.db.Statement(ctx, tx, categoryTable, fields...).
		Where("id = ?", productTagID)

	if err := stmt.First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}
