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

var notificationFields = []string{
	"id", "created_by", "creator_name", "updated_by", "title", "body",
	"published_at", "targets", "public", "created_at", "updated_at",
}

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

func (n *notification) List(
	ctx context.Context, params *ListNotificationsParams, fields ...string,
) (entity.Notifications, error) {
	var notifications entity.Notifications
	if len(fields) == 0 {
		fields = notificationFields
	}

	stmt := n.db.DB.WithContext(ctx).Table(notificationTable).Select(fields)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&notifications).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := notifications.Fill(); err != nil {
		return nil, exception.InternalError(err)
	}
	return notifications, nil
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
