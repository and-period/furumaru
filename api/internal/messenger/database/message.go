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

// var messageFields = []string{
// 	"id", "user_type", "user_id",
// 	"type", "title", "body", "link", "read",
// 	"received_at", "created_at", "updated_at",
// }

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
