package database

import (
	"context"
	"fmt"
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

func (n *notification) List(
	ctx context.Context, params *ListNotificationsParams, fields ...string,
) (entity.Notifications, error) {
	var notifications entity.Notifications

	stmt := n.db.Statement(ctx, n.db.DB, notificationTable, fields...)
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

func (n *notification) Count(ctx context.Context, params *ListNotificationsParams) (int64, error) {
	var total int64

	stmt := n.db.Count(ctx, n.db.DB, notificationTable)
	stmt = params.stmt(stmt)

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (n *notification) Get(
	ctx context.Context, notificationID string, fields ...string,
) (*entity.Notification, error) {
	notification, err := n.get(ctx, n.db.DB, notificationID, fields...)
	return notification, exception.InternalError(err)
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

func (n *notification) Update(ctx context.Context, notificationID string, params *UpdateNotificationParams) error {
	_, err := n.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := n.get(ctx, tx, notificationID); err != nil {
			return nil, err
		}

		updates := map[string]interface{}{
			"title":        params.Title,
			"body":         params.Body,
			"published_at": params.PublishedAt,
			"public":       params.Public,
			"updated_by":   params.UpdatedBy,
			"updated_at":   n.now(),
		}
		if len(params.Targets) > 0 {
			target, err := entity.Marshal(params.Targets)
			if err != nil {
				return nil, fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["targets"] = target
		}
		err := tx.WithContext(ctx).
			Table(notificationTable).
			Where("id = ?", notificationID).
			Updates(updates).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (n *notification) Delete(ctx context.Context, notificationID string) error {
	_, err := n.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := n.get(ctx, tx, notificationID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"deleted_at": n.now(),
		}
		err := tx.WithContext(ctx).
			Table(notificationTable).
			Where("id = ?", notificationID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (n *notification) get(
	ctx context.Context, tx *gorm.DB, notificationID string, fields ...string,
) (*entity.Notification, error) {
	var notification *entity.Notification

	err := n.db.Statement(ctx, tx, notificationTable, fields...).
		Where("id = ?", notificationID).
		First(&notification).Error
	if err != nil {
		return nil, err
	}
	if err := notification.Fill(); err != nil {
		return nil, err
	}
	return notification, nil
}
