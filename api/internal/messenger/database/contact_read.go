package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"gorm.io/gorm"
)

const contactReadTable = "contact_reads"

type contactRead struct {
	db  *database.Client
	now func() time.Time
}

func NewContactRead(db *database.Client) ContactRead {
	return &contactRead{
		db:  db,
		now: time.Now,
	}
}

func (c *contactRead) Create(ctx context.Context, contactRead *entity.ContactRead) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		contactRead.CreatedAt, contactRead.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(contactReadTable).Create(&contactRead).Error
		return err
	})
	return err
}
