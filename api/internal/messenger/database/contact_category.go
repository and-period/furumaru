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

func (c *contactCategory) Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error) {
	category, err := c.get(ctx, c.db.DB, categoryID, fields...)
	return category, exception.InternalError(err)
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
