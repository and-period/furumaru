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

const contactTable = "contacts"

type contact struct {
	db  *database.Client
	now func() time.Time
}

func NewContact(db *database.Client) Contact {
	return &contact{
		db:  db,
		now: jst.Now,
	}
}

func (c *contact) Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error) {
	contact, err := c.get(ctx, c.db.DB, contactID, fields...)
	return contact, exception.InternalError(err)
}

func (c *contact) get(
	ctx context.Context, tx *gorm.DB, contactID string, fields ...string,
) (*entity.Contact, error) {
	var contact *entity.Contact

	err := c.db.Statement(ctx, tx, contactTable, fields...).
		Where("id = ?", contactID).
		First(&contact).Error
	return contact, err
}
