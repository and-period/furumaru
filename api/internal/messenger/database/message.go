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

const messageTable = "messages"

var messageFields = []string{
	"id", "user_type", "user_id",
	"type", "title", "body", "link", "read",
	"received_at", "created_at", "updated_at",
}

type message struct {
	db  *database.Client
	now func() time.Time
}

func NewMessage(db *database.Client) Message {
	return &message{
		db:  db,
		now: jst.Now,
	}
}

func (m *message) List(ctx context.Context, params *ListMessagesParams, fields ...string) (entity.Messages, error) {
	var messages entity.Messages
	if len(fields) == 0 {
		fields = messageFields
	}

	stmt := m.db.DB.WithContext(ctx).Table(messageTable).Select(fields)
	stmt = params.stmt(stmt).Order("received_at DESC")
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&messages).Error
	return messages, exception.InternalError(err)
}

func (m *message) Count(ctx context.Context, params *ListMessagesParams) (int64, error) {
	var total int64

	stmt := m.db.DB.WithContext(ctx).Table(messageTable).Select("COUNT(*)")
	stmt = params.stmt(stmt)

	err := stmt.Find(&total).Error
	return total, exception.InternalError(err)
}

func (m *message) Get(ctx context.Context, messageID string, fields ...string) (*entity.Message, error) {
	message, err := m.get(ctx, m.db.DB, messageID, fields...)
	return message, exception.InternalError(err)
}

func (m *message) MultiCreate(ctx context.Context, messages entity.Messages) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := m.now()
		for i := range messages {
			messages[i].CreatedAt = now
			messages[i].UpdatedAt = now
		}

		err := tx.WithContext(ctx).Table(messageTable).Create(&messages).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (m *message) UpdateRead(ctx context.Context, messageID string) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := m.get(ctx, tx, messageID, "read")
		if err != nil {
			return nil, err
		}
		if current.Read {
			return nil, exception.ErrFailedPrecondition
		}

		params := map[string]interface{}{
			"read":       true,
			"updated_at": m.now(),
		}
		err = tx.WithContext(ctx).
			Table(messageTable).
			Where("id = ?", messageID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (m *message) get(ctx context.Context, tx *gorm.DB, messageID string, fields ...string) (*entity.Message, error) {
	var message *entity.Message
	if len(fields) == 0 {
		fields = messageFields
	}

	err := m.db.DB.WithContext(ctx).
		Table(messageTable).Select(fields).
		Where("id = ?", messageID).
		First(&message).Error
	return message, err
}
