package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const contactCategoryTable = "contact_categories"

type contactCategory struct {
	db  *mysql.Client
	now func() time.Time
}

func NewContactCategory(db *mysql.Client) database.ContactCategory {
	return &contactCategory{
		db:  db,
		now: jst.Now,
	}
}

func (c *contactCategory) List(
	ctx context.Context, params *database.ListContactCategoriesParams, fields ...string,
) (entity.ContactCategories, error) {
	var categories entity.ContactCategories

	stmt := c.db.Statement(ctx, c.db.DB, contactCategoryTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&categories).Error
	return categories, dbError(err)
}

func (c *contactCategory) MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.ContactCategories, error) {
	var categories entity.ContactCategories

	err := c.db.Statement(ctx, c.db.DB, contactCategoryTable, fields...).
		Where("id IN (?)", categoryIDs).
		Find(&categories).Error
	return categories, dbError(err)
}

func (c *contactCategory) Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error) {
	category, err := c.get(ctx, c.db.DB, categoryID, fields...)
	return category, dbError(err)
}

func (c *contactCategory) Create(ctx context.Context, category *entity.ContactCategory) error {
	now := c.now()
	category.CreatedAt, category.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(contactCategoryTable).Create(&category).Error
	return dbError(err)
}

func (c *contactCategory) get(
	ctx context.Context, tx *gorm.DB, categoryID string, fields ...string,
) (*entity.ContactCategory, error) {
	var contactCategory *entity.ContactCategory

	stmt := c.db.Statement(ctx, tx, contactCategoryTable, fields...).
		Where("id = ?", categoryID)

	if err := stmt.First(&contactCategory).Error; err != nil {
		return nil, err
	}
	return contactCategory, nil
}
