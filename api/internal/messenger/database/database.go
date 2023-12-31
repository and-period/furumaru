//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/messenger/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

var (
	ErrInvalidArgument    = &Error{err: errors.New("database: invalid argument")}
	ErrNotFound           = &Error{err: errors.New("database: not found")}
	ErrAlreadyExists      = &Error{err: errors.New("database: already exists")}
	ErrFailedPrecondition = &Error{err: errors.New("database: failed precondition")}
	ErrCanceled           = &Error{err: errors.New("database: canceled")}
	ErrDeadlineExceeded   = &Error{err: errors.New("database: deadline exceeded")}
	ErrInternal           = &Error{err: errors.New("database: internal error")}
	ErrUnknown            = &Error{err: errors.New("database: unknown")}
)

type Database struct {
	Contact         Contact
	ContactCategory ContactCategory
	ContactRead     ContactRead
	Message         Message
	MessageTemplate MessageTemplate
	Notification    Notification
	PushTemplate    PushTemplate
	ReceivedQueue   ReceivedQueue
	ReportTemplate  ReportTemplate
	Schedule        Schedule
	Thread          Thread
}

type Contact interface {
	List(ctx context.Context, params *ListContactsParams, fields ...string) (entity.Contacts, error)
	Count(ctx context.Context) (int64, error)
	Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error)
	Create(ctx context.Context, contact *entity.Contact) error
	Update(ctx context.Context, contactID string, params *UpdateContactParams) error
	Delete(ctx context.Context, contactID string) error
}

type ListContactsParams struct {
	Limit  int
	Offset int
}

type UpdateContactParams struct {
	Title       string
	CategoryID  string
	Content     string
	Username    string
	UserID      string
	Email       string
	PhoneNumber string
	Status      entity.ContactStatus
	ResponderID string
	Note        string
}

type ContactCategory interface {
	Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error)
	List(ctx context.Context, params *ListContactCategoriesParams, fields ...string) (entity.ContactCategories, error)
	MultiGet(ctx context.Context, categoryIDs []string, fields ...string) (entity.ContactCategories, error)
	Create(ctx context.Context, category *entity.ContactCategory) error
}

type ListContactCategoriesParams struct {
	Limit  int
	Offset int
}

type ContactRead interface {
	GetByContactIDAndUserID(ctx context.Context, contactID, userID string, fields ...string) (*entity.ContactRead, error)
	Create(ctx context.Context, contactRead *entity.ContactRead) error
	Update(ctx context.Context, params *UpdateContactReadParams) error
}

type UpdateContactReadParams struct {
	ContactID string
	UserID    string
	Read      bool
}

type Message interface {
	List(ctx context.Context, params *ListMessagesParams, fields ...string) (entity.Messages, error)
	Count(ctx context.Context, params *ListMessagesParams) (int64, error)
	Get(ctx context.Context, messageID string, fields ...string) (*entity.Message, error)
	MultiCreate(ctx context.Context, messages entity.Messages) error
	UpdateRead(ctx context.Context, messageID string) error
}

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

type MessageTemplate interface {
	Get(ctx context.Context, messageID entity.MessageTemplateID, fields ...string) (*entity.MessageTemplate, error)
}

type Notification interface {
	List(ctx context.Context, params *ListNotificationsParams, fields ...string) (entity.Notifications, error)
	Count(ctx context.Context, params *ListNotificationsParams) (int64, error)
	Get(ctx context.Context, notificationID string, fields ...string) (*entity.Notification, error)
	Create(ctx context.Context, notification *entity.Notification) error
	Update(ctx context.Context, notificationID string, params *UpdateNotificationParams) error
	Delete(ctx context.Context, notificationID string) error
}

type ListNotificationsParams struct {
	Limit  int
	Offset int
	Since  time.Time
	Until  time.Time
	Orders []*ListNotificationsOrder
}

type ListNotificationsOrder struct {
	Key        entity.NotificationOrderBy
	OrderByASC bool
}

type UpdateNotificationParams struct {
	Targets     []entity.NotificationTarget
	Title       string
	Body        string
	Note        string
	PublishedAt time.Time
	UpdatedBy   string
}

type PushTemplate interface {
	Get(ctx context.Context, pushID entity.PushTemplateID, fields ...string) (*entity.PushTemplate, error)
}

type ReceivedQueue interface {
	Get(ctx context.Context, queueID string, fields ...string) (*entity.ReceivedQueue, error)
	Create(ctx context.Context, queue *entity.ReceivedQueue) error
	UpdateDone(ctx context.Context, queueID string, done bool) error
}

type ReportTemplate interface {
	Get(ctx context.Context, reportID entity.ReportTemplateID, fields ...string) (*entity.ReportTemplate, error)
}

type Schedule interface {
	List(ctx context.Context, params *ListSchedulesParams, fields ...string) (entity.Schedules, error)
	Get(ctx context.Context, messageType entity.ScheduleType, messageID string, fields ...string) (*entity.Schedule, error)
	Upsert(ctx context.Context, schedule *entity.Schedule) error
	UpsertProcessing(ctx context.Context, schedule *entity.Schedule) error
	UpdateDone(ctx context.Context, messageType entity.ScheduleType, messageID string) error
	UpdateCancel(ctx context.Context, messageType entity.ScheduleType, messageID string) error
}

type ListSchedulesParams struct {
	Types    []entity.ScheduleType
	Statuses []entity.ScheduleStatus
	Since    time.Time
	Until    time.Time
}

type Thread interface {
	List(ctx context.Context, params *ListThreadsParams, fields ...string) (entity.Threads, error)
	Count(ctx context.Context, params *ListThreadsParams) (int64, error)
	Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error)
	Create(ctx context.Context, thread *entity.Thread) error
	Update(ctx context.Context, threadID string, params *UpdateThreadParams) error
	Delete(ctx context.Context, threadID string) error
}

type ListThreadsParams struct {
	ContactID string
	Limit     int
	Offset    int
}

type UpdateThreadParams struct {
	Content  string
	UserID   string
	UserType entity.ThreadUserType
}

type Error struct {
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
