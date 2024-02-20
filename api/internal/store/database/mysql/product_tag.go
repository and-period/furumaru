package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const productTagTable = "product_tags"

type productTag struct {
	db  *mysql.Client
	now func() time.Time
}

func newProductTag(db *mysql.Client) database.ProductTag {
	return &productTag{
		db:  db,
		now: jst.Now,
	}
}

type listProductTagsParams database.ListProductTagsParams

func (p listProductTagsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`name`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
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
		stmt = stmt.Order("name ASC")
	}
	return stmt
}

func (p listProductTagsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *productTag) List(
	ctx context.Context, params *database.ListProductTagsParams, fields ...string,
) (entity.ProductTags, error) {
	var tags entity.ProductTags

	p := listProductTagsParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, productTagTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&tags).Error
	return tags, dbError(err)
}

func (t *productTag) Count(ctx context.Context, params *database.ListProductTagsParams) (int64, error) {
	p := listProductTagsParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.ProductTag{}, p.stmt)
	return total, dbError(err)
}

func (t *productTag) MultiGet(
	ctx context.Context, productTagIDs []string, fields ...string,
) (entity.ProductTags, error) {
	var tags entity.ProductTags

	err := t.db.Statement(ctx, t.db.DB, productTagTable, fields...).
		Where("id IN (?)", productTagIDs).
		Find(&tags).Error
	return tags, dbError(err)
}

func (t *productTag) Get(ctx context.Context, productTagID string, fields ...string) (*entity.ProductTag, error) {
	tag, err := t.get(ctx, t.db.DB, productTagID, fields...)
	return tag, dbError(err)
}

func (t *productTag) Create(ctx context.Context, tag *entity.ProductTag) error {
	now := t.now()
	tag.CreatedAt, tag.UpdatedAt = now, now

	err := t.db.DB.WithContext(ctx).Table(productTagTable).Create(&tag).Error
	return dbError(err)
}

func (t *productTag) Update(ctx context.Context, productTagID, name string) error {
	params := map[string]interface{}{
		"name":       name,
		"updated_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).
		Table(productTagTable).
		Where("id = ?", productTagID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (t *productTag) Delete(ctx context.Context, productTagID string) error {
	stmt := t.db.DB.WithContext(ctx).
		Table(productTagTable).
		Where("id = ?", productTagID)

	err := stmt.Delete(&entity.ProductTag{}).Error
	return dbError(err)
}

func (t *productTag) get(
	ctx context.Context, tx *gorm.DB, productTagID string, fields ...string,
) (*entity.ProductTag, error) {
	var tag *entity.ProductTag

	stmt := t.db.Statement(ctx, tx, productTagTable, fields...).
		Where("id = ?", productTagID)

	if err := stmt.First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}
