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

const contactReadTable = "contact_reads"

type contactRead struct {
	db  *database.Client
	now func() time.Time
}

func NewContactRead(db *database.Client) ContactRead {
	return &contactRead{
		db:  db,
		now: jst.Now,
	}
}

func (c *contactRead) Get(ctx context.Context, contactID, userID string, fields ...string,
) (*entity.ContactRead, error) {
	contactRead, err := c.getByContactIDAndUserID(ctx, c.db.DB, contactID, userID, fields...)
	return contactRead, exception.InternalError(err)
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

func (c *contactRead) UpdateRead(ctx context.Context, params *UpdateContactReadFlagParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := c.getByContactIDAndUserID(ctx, tx, params.ContactID, params.UserID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"read":       params.Read,
			"updated_at": c.now(),
		}

		if params.UserID == "" {
			err := tx.WithContext(ctx).
				Table(contactReadTable).
				Where("contact_id = ? AND user_id IS NULL", params.ContactID).
				Updates(updates).Error
			return err
		}

		err := tx.WithContext(ctx).
			Table(contactReadTable).
			Where("contact_id = ? AND user_id = ?", params.ContactID, params.UserID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *contactRead) getByContactIDAndUserID(ctx context.Context, tx *gorm.DB, contactID, userID string, fields ...string,
) (*entity.ContactRead, error) {
	var contactRead *entity.ContactRead

	if userID == "" {
		err := c.db.Statement(ctx, tx, contactReadTable, fields...).
			Where("contact_id = ? AND user_id IS NULL", contactID).
			First(&contactRead).Error
		return contactRead, err
	}

	err := c.db.Statement(ctx, tx, contactReadTable, fields...).
		Where("contact_id = ?", contactID).
		Where("user_id = ?", userID).
		First(&contactRead).Error
	return contactRead, err
}
