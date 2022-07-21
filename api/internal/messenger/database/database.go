//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/messenger/$GOPACKAGE/$GOFILE
package database

import (
	"context"
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
	Message         Message
	MessageTemplate MessageTemplate
	Notification    Notification
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
		ReceivedQueue:   NewReceivedQueue(params.Database),
		ReportTemplate:  NewReportTemplate(params.Database),
		Schedule:        NewSchedule(params.Database),
	}
}

/**
 * interface
 */
type Contact interface {
	List(ctx context.Context, params *ListContactsParams, fields ...string) (entity.Contacts, error)
	Count(ctx context.Context, params *ListContactsParams) (int64, error)
	Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error)
	Create(ctx context.Context, contact *entity.Contact) error
	Update(ctx context.Context, contactID string, params *UpdateContactParams) error
	Delete(ctx context.Context, contactID string) error
}

type Message interface {
	MultiCreate(ctx context.Context, messages entity.Messages) error
}

type MessageTemplate interface {
	Get(ctx context.Context, messageID string, fields ...string) (*entity.MessageTemplate, error)
}

type Notification interface {
	List(ctx context.Context, params *ListNotificationsParams, fields ...string) (entity.Notifications, error)
	Get(ctx context.Context, notificationID string, fields ...string) (*entity.Notification, error)
	Create(ctx context.Context, notification *entity.Notification) error
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
type ListContactsParams struct {
	Limit  int
	Offset int
}

type UpdateContactParams struct {
	Status   entity.ContactStatus
	Priority entity.ContactPriority
	Note     string
}

type ListNotificationsParams struct {
	Limit         int
	Offset        int
	Since         time.Time
	Until         time.Time
	OnlyPublished bool
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
	return stmt
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
