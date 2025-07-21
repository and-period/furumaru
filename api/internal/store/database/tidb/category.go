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
	db  *mysql.Client
	now func() time.Time
}

func NewCategory(db *mysql.Client) database.Category {
	return &category{
		db:  db,
		now: jst.Now,
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

func (c *category) Count(
	ctx context.Context,
	params *database.ListCategoriesParams,
) (int64, error) {
	p := listCategoriesParams(*params)

	total, err := c.db.Count(ctx, c.db.DB, &entity.Category{}, p.stmt)
	return total, dbError(err)
}

func (c *category) MultiGet(
	ctx context.Context,
	categoryIDs []string,
	fields ...string,
) (entity.Categories, error) {
	var categories entity.Categories

	err := c.db.Statement(ctx, c.db.DB, categoryTable, fields...).
		Where("id IN (?)", categoryIDs).
		Find(&categories).Error
	return categories, dbError(err)
}

func (c *category) Get(
	ctx context.Context,
	categoryID string,
	fields ...string,
) (*entity.Category, error) {
	category, err := c.get(ctx, c.db.DB, categoryID, fields...)
	return category, dbError(err)
}

func (c *category) Create(ctx context.Context, category *entity.Category) error {
	now := c.now()
	category.CreatedAt, category.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(categoryTable).Create(&category).Error
	return dbError(err)
}

func (c *category) Update(ctx context.Context, categoryID, name string) error {
	params := map[string]interface{}{
		"name":       name,
		"updated_at": c.now(),
	}
	stmt := c.db.DB.WithContext(ctx).
		Table(categoryTable).
		Where("id = ?", categoryID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (c *category) Delete(ctx context.Context, categoryID string) error {
	stmt := c.db.DB.WithContext(ctx).
		Table(categoryTable).
		Where("id = ?", categoryID)

	err := stmt.Delete(&entity.Category{}).Error
	return dbError(err)
}

func (c *category) get(
	ctx context.Context, tx *gorm.DB, categoryID string, fields ...string,
) (*entity.Category, error) {
	var category *entity.Category

	stmt := c.db.Statement(ctx, tx, categoryTable, fields...).
		Where("id = ?", categoryID)

	if err := stmt.First(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}
