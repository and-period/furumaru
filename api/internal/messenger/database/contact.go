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

func (c *contact) List(ctx context.Context, params *ListContactsParams, fields ...string) (entity.Contacts, error) {
	var contacts entity.Contacts

	stmt := c.db.Statement(ctx, c.db.DB, contactTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&contacts).Error
	return contacts, exception.InternalError(err)
}

func (c *contact) Count(ctx context.Context, params *ListContactsParams) (int64, error) {
	total, err := c.db.Count(ctx, c.db.DB, &entity.Contact{}, nil)
	return total, exception.InternalError(err)
}

func (c *contact) Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error) {
	contact, err := c.get(ctx, c.db.DB, contactID, fields...)
	return contact, exception.InternalError(err)
}

func (c *contact) Create(ctx context.Context, contact *entity.Contact) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		contact.CreatedAt, contact.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(contactTable).Create(&contact).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *contact) Update(ctx context.Context, contactID string, params *UpdateContactParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := c.get(ctx, tx, contactID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":     params.Status,
			"priority":   params.Priority,
			"note":       params.Note,
			"updated_at": c.now(),
		}
		err := tx.WithContext(ctx).
			Table(contactTable).
			Where("id = ?", contactID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *contact) Delete(ctx context.Context, contactID string) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := c.get(ctx, tx, contactID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":     entity.ContactStatusDiscard,
			"updated_at": c.now(),
		}
		err := tx.WithContext(ctx).
			Table(contactTable).
			Where("id = ?", contactID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
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
