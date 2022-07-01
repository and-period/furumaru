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

const notificationTable = "notifications"

type notification struct {
	db  *database.Client
	now func() time.Time
}

func NewNotification(db *database.Client) Notification {
	return &notification{
		db:  db,
		now: jst.Now,
	}
}

func (n *notification) Create(ctx context.Context, notification *entity.Notification) error {
	_, err := n.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := n.now()
		notification.CreatedAt, notification.UpdatedAt = now, now
		err := notification.FillJSON()
		if err != nil {
			return nil, err
		}
		err = tx.WithContext(ctx).Table(notificationTable).Create(&notification).Error
		return nil, err
	})
	return exception.InternalError(err)
}
