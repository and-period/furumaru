package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const contactCategoryTable = "contact_categories"

type contactCategory struct {
	db  *database.Client
	now func() time.Time
}

func NewContactCategory(db *database.Client) ContactCategory {
	return &contactCategory{
		db:  db,
		now: jst.Now,
	}
}

func (c *contactCategory) List(
	ctx context.Context, params *ListContactCategoriesParams, fields ...string,
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
	return categories, exception.InternalError(err)
}

func (c *contactCategory) Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error) {
	category, err := c.get(ctx, c.db.DB, categoryID, fields...)
	return category, exception.InternalError(err)
}

func (c *contactCategory) Create(ctx context.Context, category *entity.ContactCategory) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		category.CreatedAt, category.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(contactCategoryTable).Create(&category).Error
		return err
	})
	return exception.InternalError(err)
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
