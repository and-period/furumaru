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

const contactTable = "contacts"

type contact struct {
	db  *mysql.Client
	now func() time.Time
}

func NewContact(db *mysql.Client) database.Contact {
	return &contact{
		db:  db,
		now: jst.Now,
	}
}

func (c *contact) List(
	ctx context.Context,
	params *database.ListContactsParams,
	fields ...string,
) (entity.Contacts, error) {
	var contacts entity.Contacts

	stmt := c.db.Statement(ctx, c.db.DB, contactTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&contacts).Error
	return contacts, dbError(err)
}

func (c *contact) Count(ctx context.Context) (int64, error) {
	total, err := c.db.Count(ctx, c.db.DB, &entity.Contact{}, nil)
	return total, dbError(err)
}

func (c *contact) Get(
	ctx context.Context,
	contactID string,
	fields ...string,
) (*entity.Contact, error) {
	contact, err := c.get(ctx, c.db.DB, contactID, fields...)
	return contact, dbError(err)
}

func (c *contact) Create(ctx context.Context, contact *entity.Contact) error {
	now := c.now()
	contact.CreatedAt, contact.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(contactTable).Create(&contact).Error
	return dbError(err)
}

func (c *contact) Update(
	ctx context.Context,
	contactID string,
	params *database.UpdateContactParams,
) error {
	updates := map[string]interface{}{
		"title":        params.Title,
		"category_id":  params.CategoryID,
		"content":      params.Content,
		"username":     params.Username,
		"user_id":      params.UserID,
		"email":        params.Email,
		"phone_number": params.PhoneNumber,
		"status":       params.Status,
		"responder_id": params.ResponderID,
		"note":         params.Note,
		"updated_at":   c.now(),
	}
	stmt := c.db.DB.WithContext(ctx).
		Table(contactTable).
		Where("id = ?", contactID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (c *contact) Delete(ctx context.Context, contactID string) error {
	params := map[string]interface{}{
		"deleted_at": c.now(),
	}
	stmt := c.db.DB.WithContext(ctx).
		Table(contactTable).
		Where("id = ?", contactID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (c *contact) get(
	ctx context.Context,
	tx *gorm.DB,
	contactID string,
	fields ...string,
) (*entity.Contact, error) {
	var contact *entity.Contact

	stmt := c.db.Statement(ctx, tx, contactTable, fields...).
		Where("id = ?", contactID)

	if err := stmt.First(&contact).Error; err != nil {
		return nil, err
	}
	return contact, nil
}
