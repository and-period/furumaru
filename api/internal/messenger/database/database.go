//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/messenger/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"gorm.io/gorm"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Contact         Contact
	ContactCategory ContactCategory
	Thread          Thread
	Message         Message
	MessageTemplate MessageTemplate
	Notification    Notification
	PushTemplate    PushTemplate
	ReceivedQueue   ReceivedQueue
	ReportTemplate  ReportTemplate
	Schedule        Schedule
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Contact:         NewContact(params.Database),
		Message:         NewMessage(params.Database),
		MessageTemplate: NewMessageTemplate(params.Database),
		Notification:    NewNotification(params.Database),
		PushTemplate:    NewPushTemplate(params.Database),
		ReceivedQueue:   NewReceivedQueue(params.Database),
		ReportTemplate:  NewReportTemplate(params.Database),
		Schedule:        NewSchedule(params.Database),
	}
}

/**
 * interface
 */

type Contact interface {
	Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error)
	Create(ctx context.Context, contact *entity.Contact) error
}

type ContactCategory interface {
	Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error)
}

type Thread interface {
	Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error)
	ListByContactID(ctx context.Context, params *ListThreadsByContactIDParams, fields ...string) (entity.Threads, error)
	Count(ctx context.Context, params *ListThreadsByContactIDParams) (int64, error)
	Create(ctx context.Context, thread *entity.Thread) error
	Update(ctx context.Context, threadID string, params *UpdateThreadParams) error
	Delete(ctx context.Context, threadID string) error
}

type Message interface {
	List(ctx context.Context, params *ListMessagesParams, fields ...string) (entity.Messages, error)
	Count(ctx context.Context, params *ListMessagesParams) (int64, error)
	Get(ctx context.Context, messageID string, fields ...string) (*entity.Message, error)
	MultiCreate(ctx context.Context, messages entity.Messages) error
	UpdateRead(ctx context.Context, messageID string) error
}

type MessageTemplate interface {
	Get(ctx context.Context, messageID string, fields ...string) (*entity.MessageTemplate, error)
}

type Notification interface {
	List(ctx context.Context, params *ListNotificationsParams, fields ...string) (entity.Notifications, error)
	Count(ctx context.Context, params *ListNotificationsParams) (int64, error)
	Get(ctx context.Context, notificationID string, fields ...string) (*entity.Notification, error)
	Create(ctx context.Context, notification *entity.Notification) error
	Update(ctx context.Context, notificationID string, params *UpdateNotificationParams) error
	Delete(ctx context.Context, notificationID string) error
}

type PushTemplate interface {
	Get(ctx context.Context, pushID string, fields ...string) (*entity.PushTemplate, error)
}

type ReceivedQueue interface {
	Get(ctx context.Context, queueID string, fields ...string) (*entity.ReceivedQueue, error)
	Create(ctx context.Context, queue *entity.ReceivedQueue) error
	UpdateDone(ctx context.Context, queueID string, done bool) error
}

type ReportTemplate interface {
	Get(ctx context.Context, reportID string, fields ...string) (*entity.ReportTemplate, error)
}

type Schedule interface {
	List(ctx context.Context, params *ListSchedulesParams, fields ...string) (entity.Schedules, error)
	UpsertProcessing(ctx context.Context, schedule *entity.Schedule) error
	UpdateDone(ctx context.Context, messageType entity.ScheduleType, messageID string) error
	UpdateCancel(ctx context.Context, messageType entity.ScheduleType, messageID string) error
}

/**
 * params
 */

type ListMessagesParams struct {
	Limit    int
	Offset   int
	UserType entity.UserType
	UserID   string
	Orders   []*ListMessagesOrder
}

type ListMessagesOrder struct {
	Key        entity.MessageOrderBy
	OrderByASC bool
}

func (p *ListMessagesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.UserType != entity.UserTypeNone {
		stmt = stmt.Where("user_type = ?", p.UserType)
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
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

type ListNotificationsParams struct {
	Limit         int
	Offset        int
	Since         time.Time
	Until         time.Time
	OnlyPublished bool
	Orders        []*ListNotificationsOrder
}

type ListNotificationsOrder struct {
	Key        entity.NotificationOrderBy
	OrderByASC bool
}

func (p *ListNotificationsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if !p.Since.IsZero() {
		stmt = stmt.Where("published_at >= ?", p.Since)
	}
	if !p.Until.IsZero() {
		stmt = stmt.Where("published_at <= ?", p.Until)
	}
	if p.OnlyPublished {
		stmt = stmt.Where("public = ?", true)
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

type UpdateNotificationParams struct {
	Title       string
	Body        string
	Targets     []entity.TargetType
	PublishedAt time.Time
	Public      bool
	UpdatedBy   string
}
type ListSchedulesParams struct {
	Types    []entity.ScheduleType
	Statuses []entity.ScheduleStatus
	Since    time.Time
	Until    time.Time
}

func (p *ListSchedulesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if len(p.Types) > 0 {
		stmt = stmt.Where("message_type IN (?)", p.Types)
	}
	if len(p.Statuses) > 0 {
		stmt = stmt.Where("status IN (?)", p.Statuses)
	}
	if !p.Since.IsZero() {
		stmt = stmt.Where("sent_at >= ?", p.Since)
	}
	if !p.Until.IsZero() {
		stmt = stmt.Where("sent_at <= ?", p.Until)
	}
	return stmt
}

type ListThreadsByContactIDParams struct {
	ContactID string
	Limit     int
	Offset    int
}

func (p *ListThreadsByContactIDParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ContactID != "" {
		stmt = stmt.Where("contact_id = ?", p.ContactID)
	}
	return stmt
}

type UpdateThreadParams struct {
	Content  string
	UserID   string
	UserType int32
}
