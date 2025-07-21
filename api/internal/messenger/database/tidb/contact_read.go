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

const contactReadTable = "contact_reads"

type contactRead struct {
	db  *mysql.Client
	now func() time.Time
}

func NewContactRead(db *mysql.Client) database.ContactRead {
	return &contactRead{
		db:  db,
		now: jst.Now,
	}
}

func (c *contactRead) GetByContactIDAndUserID(
	ctx context.Context,
	contactID, userID string,
	fields ...string,
) (*entity.ContactRead, error) {
	contactRead, err := c.getByContactIDAndUserID(ctx, c.db.DB, contactID, userID, fields...)
	return contactRead, dbError(err)
}

func (c *contactRead) Create(ctx context.Context, contactRead *entity.ContactRead) error {
	now := c.now()
	contactRead.CreatedAt, contactRead.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(contactReadTable).Create(&contactRead).Error
	return dbError(err)
}

func (c *contactRead) Update(ctx context.Context, params *database.UpdateContactReadParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
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
	return dbError(err)
}

func (c *contactRead) getByContactIDAndUserID(
	ctx context.Context,
	tx *gorm.DB,
	contactID, userID string,
	fields ...string,
) (*entity.ContactRead, error) {
	var contactRead *entity.ContactRead

	stmt := c.db.Statement(ctx, tx, contactReadTable, fields...).Where("contact_id = ?", contactID)
	if userID == "" {
		stmt = stmt.Where("user_id IS NULL")
	} else {
		stmt = stmt.Where("user_id = ?", userID)
	}

	if err := stmt.First(&contactRead).Error; err != nil {
		return nil, err
	}
	return contactRead, nil
}
