package messenger

import (
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

/**
 * Concact - お問い合わせ
 */
type ListContactsInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type GetContactInput struct {
	ContactID string `validate:"required"`
}

type CreateContactInput struct {
	Title       string `validate:"required,max=128"`
	Content     string `validate:"required,max=2000"`
	Username    string `validate:"required,max=128"`
	UserID      string `validate:""`
	CategoryID  string `validate:"required,max=128"`
	Email       string `validate:"required,max=256,email"`
	PhoneNumber string `validate:"required,e164"`
	ResponderID string `validate:""`
	Note        string `validate:"max=2000"`
}

type UpdateContactInput struct {
	ContactID   string               `validate:"required"`
	Title       string               `validate:"required,max=128"`
	Content     string               `validate:"required,max=2000"`
	Username    string               `validate:"required,max=128"`
	UserID      string               `validate:""`
	CategoryID  string               `validate:"required,max=128"`
	Email       string               `validate:"required,max=256,email"`
	PhoneNumber string               `validate:"required,e164"`
	Status      entity.ContactStatus `validate:"required"`
	ResponderID string               `validate:""`
	Note        string               `validate:"max=2000"`
}

type DeleteContactInput struct {
	ContactID string `validate:"required"`
}

/**
 * ContactCategory - お問い合わせ種別
 */
type ListContactCategoriesInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type MultiGetContactCategoriesInput struct {
	CategoryIDs []string `validate:"opitempty,dive,required"`
}

type GetContactCategoryInput struct {
	CategoryID string `validate:"required"`
}

/**
 * ContactRead - お問い合わせ既読管理
 */
type CreateContactReadInput struct {
	ContactID string                 `validate:"required"`
	UserID    string                 `validate:""`
	UserType  entity.ContactUserType `validate:"required"`
}

/**
 * Message - メッセージ
 */
type ListMessagesInput struct {
	UserType entity.UserType      `validate:"required,oneof=1 2"`
	UserID   string               `validate:"required"`
	Limit    int64                `validate:"required,max=200"`
	Offset   int64                `validate:"min=0"`
	Orders   []*ListMessagesOrder `validate:"dive,required"`
}

type ListMessagesOrder struct {
	Key        entity.MessageOrderBy `validate:"required"`
	OrderByASC bool                  `validate:""`
}

type GetMessageInput struct {
	MessageID string          `validate:"required"`
	UserType  entity.UserType `validate:"omitempty,oneof=1 2"`
	UserID    string          `validate:""`
}

/**
 * Notification - お知らせ
 */
type ListNotificationsOrderKey int32

const (
	ListNotificationsOrderByTitle ListNotificationsOrderKey = iota + 1
	ListNotificationsOrderByPublishedAt
)

type ListNotificationsInput struct {
	Limit  int64                     `validate:"required,max=200"`
	Offset int64                     `validate:"min=0"`
	Since  time.Time                 `validate:""`
	Until  time.Time                 `validate:""`
	Orders []*ListNotificationsOrder `validate:"dive,required"`
}

type ListNotificationsOrder struct {
	Key        ListNotificationsOrderKey `validate:"required"`
	OrderByASC bool                      `validate:""`
}

type GetNotificationInput struct {
	NotificationID string `validate:"required"`
}

type CreateNotificationInput struct {
	Type        entity.NotificationType     `validate:"required"`
	Title       string                      `validate:"max=128"`
	Body        string                      `validate:"required,max=2000"`
	Note        string                      `validate:"max=2000"`
	Targets     []entity.NotificationTarget `validate:"min=1,max=4,unique,dive,required"`
	PublishedAt time.Time                   `validate:"required"`
	CreatedBy   string                      `validate:"required"`
	PromotionID string                      `validate:""`
}

type UpdateNotificationInput struct {
	NotificationID string                      `validate:"required"`
	Title          string                      `validate:"max=128"`
	Body           string                      `validate:"required,max=2000"`
	Note           string                      `validate:"max=2000"`
	Targets        []entity.NotificationTarget `validate:"min=1,max=4,unique,dive,required"`
	PublishedAt    time.Time                   `validate:"required"`
	UpdatedBy      string                      `validete:"required"`
}

type DeleteNotificationInput struct {
	NotificationID string `validate:"required"`
}

/**
 * Notify - 通知関連(共通)
 */
type NotifyNotificationInput struct {
	NotificationID string `validate:"required"`
}

/**
 * NotifyAdmin - 通知関連(管理者宛)
 */
type NotifyRegisterAdminInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type NotifyResetAdminPasswordInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

/**
 * NotifyUser - 通知関連(利用者宛)
 */
type NotifyStartLiveInput struct {
	ScheduleID string `validate:"required"`
}

type NotifyOrderCapturedInput struct {
	OrderID string `validate:"required"`
}

type NotifyOrderShippedInput struct {
	OrderID string `validate:"required"`
}

/**
 * ReserveNotification - 通知予約関連
 */
type ReserveStartLiveInput struct {
	ScheduleID string `validate:"required"`
}

type ReserveNotificationInput struct {
	NotificationID string `validate:"required"`
}

/**
 * Thread - お問い合わせ会話履歴
 */
type ListThreadsInput struct {
	ContactID string `validate:"required"`
	UserID    string `validate:""`
	Limit     int64  `validate:"required,max=200"`
	Offset    int64  `validate:"min=0"`
}

type GetThreadInput struct {
	ThreadID string `validate:"required"`
}

type CreateThreadInput struct {
	ContactID string                `validate:"required"`
	UserID    string                `validate:""`
	UserType  entity.ThreadUserType `validate:"required"`
	Content   string                `validate:"required,max=2000"`
}

type UpdateThreadInput struct {
	ThreadID string                `validate:"required"`
	Content  string                `validate:"required,max=2000"`
	UserID   string                `validate:""`
	UserType entity.ThreadUserType `validate:"required"`
}

type DeleteThreadInput struct {
	ThreadID string `validate:"required"`
}
