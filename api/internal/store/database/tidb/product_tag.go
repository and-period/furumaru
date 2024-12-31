package tidb

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
	database.ProductTag
	db  *mysql.Client
	now func() time.Time
}

func NewProductTag(db *mysql.Client, mysql database.ProductTag) database.ProductTag {
	return &productTag{
		ProductTag: mysql,
		db:         db,
		now:        jst.Now,
	}
}

type listProductTagsParams database.ListProductTagsParams

func (p listProductTagsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
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
