package tidb

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const userNotificationTable = "user_notifications"

type userNotification struct {
	db  *mysql.Client
	now func() time.Time
}

func NewUserNotification(db *mysql.Client) database.UserNotification {
	return &userNotification{
		db:  db,
		now: jst.Now,
	}
}

func (n *userNotification) MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.UserNotifications, error) {
	var notifications entity.UserNotifications

	stmt := n.db.Statement(ctx, n.db.DB, userNotificationTable, fields...).Where("user_id IN (?)", userIDs)

	err := stmt.Find(&notifications).Error
	return notifications, dbError(err)
}

func (n *userNotification) Get(ctx context.Context, userID string, fields ...string) (*entity.UserNotification, error) {
	notification, err := n.get(ctx, n.db.DB, userID, fields...)
	return notification, dbError(err)
}

func (n *userNotification) Upsert(ctx context.Context, notification *entity.UserNotification) error {
	now := n.now()
	notification.CreatedAt, notification.UpdatedAt = now, now

	updates := map[string]interface{}{
		"disabled":   notification.Disabled,
		"updated_at": now,
	}
	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(updates),
	}
	err := n.db.DB.WithContext(ctx).Clauses(clauses).Create(&notification).Error
	return dbError(err)
}

func (n *userNotification) get(ctx context.Context, tx *gorm.DB, userID string, fields ...string) (*entity.UserNotification, error) {
	var notification *entity.UserNotification

	stmt := n.db.Statement(ctx, tx, userNotificationTable, fields...).Where("user_id = ?", userID)

	err := stmt.First(&notification).Error
	if err == nil {
		return notification, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	now := n.now()
	notification = entity.NewUserNotification(userID)
	notification.CreatedAt, notification.UpdatedAt = now, now

	if err := tx.WithContext(ctx).Table(userNotificationTable).Create(&notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}
