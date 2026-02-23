package tidb

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const notificationTable = "notifications"

type notification struct {
	db  *mysql.Client
	now func() time.Time
}

func NewNotification(db *mysql.Client) database.Notification {
	return &notification{
		db:  db,
		now: jst.Now,
	}
}

type listNotificationsParams database.ListNotificationsParams

func (p listNotificationsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if !p.Since.IsZero() {
		stmt = stmt.Where("published_at >= ?", p.Since)
	}
	if !p.Until.IsZero() {
		stmt = stmt.Where("published_at <= ?", p.Until)
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

func (p listNotificationsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (n *notification) List(
	ctx context.Context, params *database.ListNotificationsParams, fields ...string,
) (entity.Notifications, error) {
	var internal internalNotifications

	p := listNotificationsParams(*params)

	stmt := n.db.Statement(ctx, n.db.DB, notificationTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	notifications := internal.entities()

	notifications.Fill(n.now())
	return notifications, nil
}

func (n *notification) Count(ctx context.Context, params *database.ListNotificationsParams) (int64, error) {
	p := listNotificationsParams(*params)

	total, err := n.db.Count(ctx, n.db.DB, &entity.Notification{}, p.stmt)
	return total, dbError(err)
}

func (n *notification) Get(
	ctx context.Context, notificationID string, fields ...string,
) (*entity.Notification, error) {
	notification, err := n.get(ctx, n.db.DB, notificationID, fields...)
	return notification, dbError(err)
}

func (n *notification) Create(ctx context.Context, notification *entity.Notification) error {
	now := n.now()
	notification.CreatedAt, notification.UpdatedAt = now, now

	internal := newInternalNotification(notification)

	err := n.db.DB.WithContext(ctx).Table(notificationTable).Create(&internal).Error
	return dbError(err)
}

func (n *notification) Update(ctx context.Context, notificationID string, params *database.UpdateNotificationParams) error {
	err := n.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := n.get(ctx, tx, notificationID)
		if err != nil {
			return err
		}
		if n.now().After(current.PublishedAt) {
			return database.ErrFailedPrecondition
		}

		updates := map[string]interface{}{
			"title":        params.Title,
			"body":         params.Body,
			"note":         params.Note,
			"published_at": params.PublishedAt,
			"updated_by":   params.UpdatedBy,
			"updated_at":   n.now(),
		}
		if len(params.Targets) > 0 {
			targetsVal, err := mysql.NewJSONColumn(params.Targets).Value()
			if err != nil {
				return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
			}
			updates["targets"] = targetsVal
		}
		err = tx.WithContext(ctx).
			Table(notificationTable).
			Where("id = ?", notificationID).
			Updates(updates).Error
		return err
	})
	return dbError(err)
}

func (n *notification) Delete(ctx context.Context, notificationID string) error {
	params := map[string]interface{}{
		"deleted_at": n.now(),
	}
	stmt := n.db.DB.WithContext(ctx).
		Table(notificationTable).
		Where("id = ?", notificationID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (n *notification) get(
	ctx context.Context, tx *gorm.DB, notificationID string, fields ...string,
) (*entity.Notification, error) {
	var internal *internalNotification

	stmt := n.db.Statement(ctx, tx, notificationTable, fields...).
		Where("id = ?", notificationID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, err
	}

	notification := internal.entity()
	notification.Fill(n.now())
	return notification, nil
}

type internalNotification struct {
	entity.Notification `gorm:"embedded"`
	TargetsJSON         mysql.JSONColumn[[]entity.NotificationTarget] `gorm:"default:null;column:targets"` // お知らせ通知先一覧(JSON)
}

type internalNotifications []*internalNotification

func newInternalNotification(notification *entity.Notification) *internalNotification {
	return &internalNotification{
		Notification: *notification,
		TargetsJSON:  mysql.NewJSONColumn(notification.Targets),
	}
}

func (n *internalNotification) entity() *entity.Notification {
	n.Notification.Targets = n.TargetsJSON.Val
	return &n.Notification
}

func (ns internalNotifications) entities() entity.Notifications {
	res := make(entity.Notifications, len(ns))
	for i := range ns {
		res[i] = ns[i].entity()
	}
	return res
}
