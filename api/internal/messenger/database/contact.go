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

var contactFields = []string{
	"id", "title", "content", "username", "email", "phone_number",
	"status", "priority", "note", "created_at", "updated_at",
}

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
	if len(fields) == 0 {
		fields = contactFields
	}

	stmt := c.db.DB.WithContext(ctx).Table(contactTable).Select(fields)
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
	var total int64

	stmt := c.db.DB.WithContext(ctx).Table(contactTable).Select("COUNT(*)")

	err := stmt.Find(&total).Error
	return total, exception.InternalError(err)
}

func (c *contact) Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error) {
	contact, err := c.get(ctx, c.db.DB, contactID, fields...)
	return contact, exception.InternalError(err)
}

func (c *contact) Create(ctx context.Context, contact *entity.Contact) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := c.now()
		contact.CreatedAt, contact.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(contactTable).Create(&contact).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *contact) Update(ctx context.Context, contactID string, params *UpdateContactParams) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, contactID); err != nil {
			return nil, err
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
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *contact) Delete(ctx context.Context, contactID string) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, contactID); err != nil {
			return nil, err
		}

		updates := map[string]interface{}{
			"status":     entity.ContactStatusDiscard,
			"updated_at": c.now(),
		}
		err := tx.WithContext(ctx).
			Table(contactTable).
			Where("id = ?", contactID).
			Updates(updates).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *contact) get(
	ctx context.Context, tx *gorm.DB, contactID string, fields ...string,
) (*entity.Contact, error) {
	var contact *entity.Contact
	if len(fields) == 0 {
		fields = contactFields
	}

	err := tx.WithContext(ctx).
		Table(contactTable).Select(fields).
		Where("id = ?", contactID).
		First(&contact).Error
	return contact, err
}
