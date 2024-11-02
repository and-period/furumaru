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

const categoryTable = "categories"

type category struct {
	database.Category
	db  *mysql.Client
	now func() time.Time
}

func newCategory(db *mysql.Client, mysql database.Category) database.Category {
	return &category{
		Category: mysql,
		db:       db,
		now:      jst.Now,
	}
}

type listCategoriesParams database.ListCategoriesParams

func (p listCategoriesParams) stmt(stmt *gorm.DB) *gorm.DB {
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

func (p listCategoriesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (c *category) List(
	ctx context.Context, params *database.ListCategoriesParams, fields ...string,
) (entity.Categories, error) {
	var categories entity.Categories

	p := listCategoriesParams(*params)

	stmt := c.db.Statement(ctx, c.db.DB, categoryTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&categories).Error
	return categories, dbError(err)
}

func (c *category) Count(ctx context.Context, params *database.ListCategoriesParams) (int64, error) {
	p := listCategoriesParams(*params)

	total, err := c.db.Count(ctx, c.db.DB, &entity.Category{}, p.stmt)
	return total, dbError(err)
}
